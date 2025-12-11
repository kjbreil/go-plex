package plex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"nhooyr.io/websocket"
)

// TimelineEntry ...
type TimelineEntry struct {
	Identifier    string `json:"identifier"`
	ItemID        int64  `json:"itemID"`
	MetadataState string `json:"metadataState"`
	SectionID     string `json:"sectionID"`
	State         int64  `json:"state"`
	Title         string `json:"title"`
	Type          int64  `json:"type"`
	UpdatedAt     int64  `json:"updatedAt"`
}

// ActivityNotification ...
type ActivityNotification struct {
	Activity struct {
		Cancellable bool   `json:"cancellable"`
		Progress    int64  `json:"progress"`
		Subtitle    string `json:"subtitle"`
		Title       string `json:"title"`
		Type        string `json:"type"`
		UserID      int64  `json:"userID"`
		UUID        string `json:"uuid"`
	} `json:"Activity"`
	Event string `json:"event"`
	UUID  string `json:"uuid"`
}

// StatusNotification ...
type StatusNotification struct {
	Description      string `json:"description"`
	NotificationName string `json:"notificationName"`
	Title            string `json:"title"`
}

// PlaySessionStateNotification ...
type PlaySessionStateNotification struct {
	GUID             string `json:"guid"`
	Key              string `json:"key"`
	PlayQueueItemID  int64  `json:"playQueueItemID"`
	RatingKey        string `json:"ratingKey"`
	SessionKey       string `json:"sessionKey"`
	State            string `json:"state"`
	URL              string `json:"url"`
	ViewOffset       int64  `json:"viewOffset"`
	TranscodeSession string `json:"transcodeSession"`
}

// ReachabilityNotification ...
type ReachabilityNotification struct {
	Reachability bool `json:"reachability"`
}

// BackgroundProcessingQueueEventNotification ...
type BackgroundProcessingQueueEventNotification struct {
	Event   string `json:"event"`
	QueueID int64  `json:"queueID"`
}

// TranscodeSession ...
type TranscodeSession struct {
	AudioChannels        int64   `json:"audioChannels"`
	AudioCodec           string  `json:"audioCodec"`
	AudioDecision        string  `json:"audioDecision"`
	Complete             bool    `json:"complete"`
	Container            string  `json:"container"`
	Context              string  `json:"context"`
	Duration             int64   `json:"duration"`
	Key                  string  `json:"key"`
	Progress             float64 `json:"progress"`
	Protocol             string  `json:"protocol"`
	Remaining            int64   `json:"remaining"`
	SourceAudioCodec     string  `json:"sourceAudioCodec"`
	SourceVideoCodec     string  `json:"sourceVideoCodec"`
	Speed                float64 `json:"speed"`
	Throttled            bool    `json:"throttled"`
	TranscodeHwRequested bool    `json:"transcodeHwRequested"`
	VideoCodec           string  `json:"videoCodec"`
	VideoDecision        string  `json:"videoDecision"`
}

// Setting ...
type Setting struct {
	Advanced bool   `json:"advanced"`
	Default  string `json:"default"`
	Group    string `json:"group"`
	Hidden   bool   `json:"hidden"`
	ID       string `json:"id"`
	Label    string `json:"label"`
	Summary  string `json:"summary"`
	Type     string `json:"type"`
	Value    int64  `json:"value"`
}

// NotificationContainer read pms notifications
type NotificationContainer struct {
	TimelineEntry []TimelineEntry `json:"TimelineEntry"`

	ActivityNotification []ActivityNotification `json:"ActivityNotification"`

	StatusNotification []StatusNotification `json:"StatusNotification"`

	PlaySessionStateNotification []PlaySessionStateNotification `json:"PlaySessionStateNotification"`

	ReachabilityNotification []ReachabilityNotification `json:"ReachabilityNotification"`

	BackgroundProcessingQueueEventNotification []BackgroundProcessingQueueEventNotification `json:"BackgroundProcessingQueueEventNotification"`

	TranscodeSession []TranscodeSession `json:"TranscodeSession"`

	Setting []Setting `json:"Setting"`

	Size int64 `json:"size"`
	// Type can be one of:
	// playing,
	// reachability,
	// transcode.end,
	// preference,
	// update.statechange,
	// activity,
	// backgroundProcessingQueue,
	// transcodeSession.update
	// transcodeSession.end
	Type string `json:"type"`
}

// WebsocketNotification websocket payload of notifications from a plex media server
type WebsocketNotification struct {
	NotificationContainer `json:"NotificationContainer"`
}

// NotificationEvents hold callbacks that correspond to notifications
type NotificationEvents struct {
	events map[string]EventHandler
}

// NewNotificationEvents initializes the event callbacks
func NewNotificationEvents() *NotificationEvents {
	return &NotificationEvents{
		events: make(map[string]EventHandler),
	}
}

type EventHandler func(n NotificationContainer)

// OnPlaying shows state information (resume, stop, pause) on a user consuming media in plex
func (e *NotificationEvents) OnPlaying(fn func(n NotificationContainer)) {
	e.events["playing"] = fn
}

// OnTranscodeUpdate shows transcode information when a transcoding stream changes parameters
func (e *NotificationEvents) OnTranscodeUpdate(fn func(n NotificationContainer)) {
	e.events["transcodeSession.update"] = fn
}

// activity
func (e *NotificationEvents) OnActivity(fn func(n NotificationContainer)) {
	e.events["activity"] = fn
}

// update.statechange
func (e *NotificationEvents) OnUpdateStateChange(fn func(n NotificationContainer)) {
	e.events["update.statechange"] = fn
}

// SubscribeToNotifications connects to your server via websockets listening for events
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

			var notif WebsocketNotification

			if err = json.Unmarshal(message, &notif); err != nil {
				p.logger.Error("websocket convert message to json failed", "err", err.Error())
				continue
			}

			if fn, ok := p.Websocket.events[notif.Type]; ok {
				fn(notif.NotificationContainer)

			} else {
				p.logger.Error("no event handler for type", "type", notif.Type, "message", notif.NotificationContainer)
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
