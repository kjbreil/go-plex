package plex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/kjbreil/go-plex/internal/plex/notification"
	"nhooyr.io/websocket"
)

// NotificationContainer is an alias for the internal notification container type.
type NotificationContainer = notification.Container

// NotificationEvents hold callbacks that correspond to notifications.
type NotificationEvents struct {
	events map[string]EventHandler
}

// NewNotificationEvents initializes the event callbacks.
func NewNotificationEvents() *NotificationEvents {
	return &NotificationEvents{
		events: make(map[string]EventHandler),
	}
}

type EventHandler func(n NotificationContainer)

// OnPlaying shows state information (resume, stop, pause) on a user consuming media in plex.
func (e *NotificationEvents) OnPlaying(fn func(n NotificationContainer)) {
	e.events["playing"] = fn
}

// OnTranscodeUpdate shows transcode information when a transcoding stream changes parameters.
func (e *NotificationEvents) OnTranscodeUpdate(fn func(n NotificationContainer)) {
	e.events["transcodeSession.update"] = fn
}

// OnActivity handles activity notifications.
func (e *NotificationEvents) OnActivity(fn func(n NotificationContainer)) {
	e.events["activity"] = fn
}

// OnUpdateStateChange handles update state change notifications.
func (e *NotificationEvents) OnUpdateStateChange(fn func(n NotificationContainer)) {
	e.events["update.statechange"] = fn
}

// SubscribeToNotifications connects to your server via websockets listening for events.
func (p *Plex) SubscribeToNotifications() {
	websocketURL := url.URL{Scheme: "ws", Host: p.url.Host, Path: "/:/websockets/notifications"}

	dialOpts := &websocket.DialOptions{
		HTTPHeader: make(http.Header),
	}
	dialOpts.HTTPHeader.Set("X-Plex-Token", p.token)

	c, _, err := websocket.Dial(p.ctx, websocketURL.String(), dialOpts)

	if err != nil {
		p.logger.Error("could not dial websocket", "err", err.Error())
		return
	}

	go func() {
		p.wg.Add(1)

		defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
			p.wg.Done()
			_ = c.CloseNow()
		}(c, websocket.StatusNormalClosure, "")

		for {
			messageType, message, err := c.Read(p.ctx)

			if err != nil {
				if p.ctx.Err() == nil {
					// Probably need to reconnect
					p.logger.Error("could not read websocket", "err", err.Error())
				} else {
					return
				}
			}

			switch messageType {
			case websocket.MessageText:
			case websocket.MessageBinary:
			}

			var notif notification.WebsocketNotification

			if err = json.Unmarshal(message, &notif); err != nil {
				p.logger.Error("websocket convert message to json failed", "err", err.Error())
				continue
			}

			if fn, ok := p.Websocket.events[notif.Type]; ok {
				fn(notif.Container)
			} else {
				p.logger.Error("no event handler for type", "type", notif.Type, "message", notif.Container)
			}
		}
	}()

	go func() {
		p.wg.Add(1)
		defer p.wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var err error

		for {
			select {
			case t := <-ticker.C:
				err = c.Write(p.ctx, websocket.MessageText, []byte(t.String()))

				if err != nil {
					p.logger.Error("failed to write to websocket", "err", err.Error())
				}
			case <-p.ctx.Done():
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				err = c.Close(websocket.StatusNormalClosure, "")

				if err != nil {
					p.logger.Error("failed to close websocket", "err", err.Error())
				}
				return
			}
		}
	}()
}
