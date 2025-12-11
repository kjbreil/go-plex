package plex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Webhook holds the actions for each webhook events
type Webhook struct {
	events map[string]func(w WebhookEvent)
	port   int
	ips    []net.IP
}

func NewWebhook(port int, ips ...net.IP) *Webhook {
	return &Webhook{
		port: port,
		ips:  ips,
		events: map[string]func(w WebhookEvent){
			"media.play":     func(w WebhookEvent) {},
			"media.pause":    func(w WebhookEvent) {},
			"media.resume":   func(w WebhookEvent) {},
			"media.stop":     func(w WebhookEvent) {},
			"media.scrobble": func(w WebhookEvent) {},
			"media.rate":     func(w WebhookEvent) {},
		},
	}
}

func (p *Plex) ServeWebhook() {
	for _, ip := range p.Webhook.ips {

		hookUrl := fmt.Sprintf("http://%s:%d/", ip.String(), p.Webhook.port)

		hooks, err := p.getWebhooks()
		if err != nil {
			panic(err)
		}

		var exists bool
		for _, hook := range hooks {
			if hook == hookUrl {
				exists = true
			}
		}
		if !exists {
			err = p.addWebhook(hookUrl)
			if err != nil {
				panic(err)
			}
		}

		http.HandleFunc("/", p.Webhook.handler)

		go func() {
			err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip.String(), p.Webhook.port), nil)
			if err != nil {
				panic(err)
			}
		}()

	}
}

// Handler listens for plex webhooks and executes the corresponding function
func (wh *Webhook) handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(0); err != nil {
		fmt.Printf("can not read form: %v", err)
		return
	}

	var hookEvent WebhookEvent

	payload, hasPayload := r.MultipartForm.Value["payload"]

	if hasPayload {

		if err := json.Unmarshal([]byte(payload[0]), &hookEvent); err != nil {
			fmt.Printf("can not parse json: %v", err)
			return
		}

		fn, ok := wh.events[hookEvent.Event]

		if !ok {
			fmt.Printf("unknown event name: %v\n", hookEvent.Event)
			return
		}

		fn(hookEvent)
	}
}

// newWebhookEvent attaches a function to each webhook event
func (wh *Webhook) newWebhookEvent(eventName string, onEvent func(w WebhookEvent)) error {
	switch eventName {
	case "media.play":
	case "media.pause":
	case "media.resume":
	case "media.stop":
	case "media.scrobble":
	case "media.rate":

	default:
		return errors.New("invalid event name")
	}

	wh.events[eventName] = onEvent

	return nil
}

// NewWebhook inits and returns a webhook event

// OnPlay executes when the webhook receives a play event
func (wh *Webhook) OnPlay(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.play", fn)
}

// OnPause executes when the webhook receives a pause event
func (wh *Webhook) OnPause(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.pause", fn)
}

// OnResume executes when the webhook receives a resume event
func (wh *Webhook) OnResume(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.resume", fn)
}

// OnStop executes when the webhook receives a stop event
func (wh *Webhook) OnStop(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.stop", fn)
}

// OnScrobble executes when the webhook receives a scrobble event
func (wh *Webhook) OnScrobble(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.scrobble", fn)
}

// OnRate executes when the webhook receives a rate event
func (wh *Webhook) OnRate(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.rate", fn)
}
