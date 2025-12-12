package plex

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	webhook "github.com/kjbreil/go-plex/internal/plex/webhook"
)

// WebhookEvent is an alias for the internal webhook event type.
type WebhookEvent = webhook.Event

// Webhook holds the actions for each webhook events.
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
			"media.play":     func(_ WebhookEvent) {},
			"media.pause":    func(_ WebhookEvent) {},
			"media.resume":   func(_ WebhookEvent) {},
			"media.stop":     func(_ WebhookEvent) {},
			"media.scrobble": func(_ WebhookEvent) {},
			"media.rate":     func(_ WebhookEvent) {},
		},
	}
}

func (p *Plex) ServeWebhook() {
	for _, ip := range p.Webhook.ips {
		hookURL := "http://" + net.JoinHostPort(ip.String(), strconv.Itoa(p.Webhook.port)) + "/"

		hooks, err := p.getWebhooks()
		if err != nil {
			panic(err)
		}

		var exists bool
		for _, hook := range hooks {
			if hook == hookURL {
				exists = true
			}
		}
		if !exists {
			err = p.addWebhook(hookURL)
			if err != nil {
				panic(err)
			}
		}

		http.HandleFunc("/", p.Webhook.handler)

		go func() {
			server := &http.Server{
				Addr:              net.JoinHostPort(ip.String(), strconv.Itoa(p.Webhook.port)),
				ReadHeaderTimeout: 10 * time.Second,
			}
			if serveErr := server.ListenAndServe(); serveErr != nil {
				p.logger.Error("webhook server error", "err", serveErr)
			}
		}()
	}
}

// Handler listens for plex webhooks and executes the corresponding function.
func (wh *Webhook) handler(_ http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(0); err != nil {
		return
	}

	var hookEvent WebhookEvent

	payload, hasPayload := r.MultipartForm.Value["payload"]

	if hasPayload {
		if err := json.Unmarshal([]byte(payload[0]), &hookEvent); err != nil {
			return
		}

		fn, ok := wh.events[hookEvent.Event]

		if !ok {
			return
		}

		fn(hookEvent)
	}
}

// newWebhookEvent attaches a function to each webhook event.
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

// OnPlay executes when the webhook receives a play event.
func (wh *Webhook) OnPlay(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.play", fn)
}

// OnPause executes when the webhook receives a pause event.
func (wh *Webhook) OnPause(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.pause", fn)
}

// OnResume executes when the webhook receives a resume event.
func (wh *Webhook) OnResume(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.resume", fn)
}

// OnStop executes when the webhook receives a stop event.
func (wh *Webhook) OnStop(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.stop", fn)
}

// OnScrobble executes when the webhook receives a scrobble event.
func (wh *Webhook) OnScrobble(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.scrobble", fn)
}

// OnRate executes when the webhook receives a rate event.
func (wh *Webhook) OnRate(fn func(w WebhookEvent)) error {
	return wh.newWebhookEvent("media.rate", fn)
}

// Webhook setup functions

type webhookHooks struct {
	URL string `json:"url"`
}

func (p *Plex) getWebhooks() ([]string, error) {
	var webhooks []string

	endpoint := "/api/v2/user/webhooks/"

	resp, err := getHost[[]webhookHooks](p, PlexURL, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range resp {
		webhooks = append(webhooks, h.URL)
	}

	return webhooks, nil
}

func (p *Plex) addWebhook(webhookURL string) error {
	// get current webhooks and append ours to it
	currentWebhooks, err := p.getWebhooks()

	if err != nil {
		return err
	}

	currentWebhooks = append(currentWebhooks, webhookURL)

	return p.setWebhooks(currentWebhooks)
}

func (p *Plex) removeWebhooks() {
	if p.Webhook == nil {
		return
	}
	for _, ip := range p.Webhook.ips {
		hookURL := "http://" + net.JoinHostPort(ip.String(), strconv.Itoa(p.Webhook.port)) + "/"

		err := p.removeWebhook(hookURL)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plex) removeWebhook(webhookURL string) error {
	currentWebhooks, err := p.getWebhooks()

	if err != nil {
		return err
	}

	for i, h := range currentWebhooks {
		if h == webhookURL {
			currentWebhooks = append(currentWebhooks[:i], currentWebhooks[i+1:]...)
			break
		}
	}

	return p.setWebhooks(currentWebhooks)
}

// SetWebhooks will set your webhooks to whatever you pass as an argument
// webhooks with a length of 0 will remove all webhooks.
func (p *Plex) setWebhooks(webhooks []string) error {
	endpoint := "/api/v2/user/webhooks"

	body := url.Values{}

	if len(webhooks) == 0 {
		body.Add("urls[]", "")
	}

	for _, hook := range webhooks {
		body.Add("urls[]", hook)
	}

	err := postHost(p, PlexURL, endpoint, []byte(body.Encode()))
	if err != nil {
		return err
	}
	return nil
}
