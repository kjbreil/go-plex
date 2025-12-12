package plex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"
	"github.com/kjbreil/go-plex/internal/plex/notification"
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
		HTTPHeader:           make(http.Header),
		HTTPClient:           nil,
		Host:                 "",
		Subprotocols:         nil,
		CompressionMode:      0,
		CompressionThreshold: 0,
		OnPingReceived:       nil,
		OnPongReceived:       nil,
	}
	dialOpts.HTTPHeader.Set("X-Plex-Token", p.token)

	c, resp, dialErr := websocket.Dial(p.ctx, websocketURL.String(), dialOpts)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if dialErr != nil {
		p.logger.Error("could not dial websocket", "err", dialErr.Error())
		return
	}

	p.wg.Add(1)
	go func() {
		defer func(_ *websocket.Conn, _ websocket.StatusCode, _ string) {
			p.wg.Done()
			_ = c.CloseNow()
		}(c, websocket.StatusNormalClosure, "")

		for {
			messageType, message, readErr := c.Read(p.ctx)

			if readErr != nil {
				if p.ctx.Err() == nil {
					// Probably need to reconnect
					p.logger.Error("could not read websocket", "err", readErr.Error())
				} else {
					return
				}
			}

			switch messageType {
			case websocket.MessageText:
			case websocket.MessageBinary:
			}

			var notif notification.WebsocketNotification

			if unmarshalErr := json.Unmarshal(message, &notif); unmarshalErr != nil {
				p.logger.Error("websocket convert message to json failed", "err", unmarshalErr.Error())
				continue
			}

			if fn, ok := p.Websocket.events[notif.Type]; ok {
				fn(notif.Container)
			} else {
				p.logger.Error("no event handler for type", "type", notif.Type, "message", notif.Container)
			}
		}
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				writeErr := c.Write(p.ctx, websocket.MessageText, []byte(t.String()))

				if writeErr != nil {
					p.logger.Error("failed to write to websocket", "err", writeErr.Error())
				}
			case <-p.ctx.Done():
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				closeErr := c.Close(websocket.StatusNormalClosure, "")

				if closeErr != nil {
					p.logger.Error("failed to close websocket", "err", closeErr.Error())
				}
				return
			}
		}
	}()
}
