package plex

import (
	"fmt"
	"net/url"
)

func (p *Plex) getWebhooks() ([]string, error) {
	type Hooks struct {
		URL string `json:"url"`
	}

	var webhooks []string

	endpoint := "/api/v2/user/webhooks/"

	resp, err := GetHost[[]Hooks](p, PlexURL, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range resp {
		webhooks = append(webhooks, h.URL)
	}

	return webhooks, nil
}

func (p *Plex) addWebhook(webhook string) error {
	// get current webhooks and append ours to it
	currentWebhooks, err := p.getWebhooks()

	if err != nil {
		return err
	}

	currentWebhooks = append(currentWebhooks, webhook)

	return p.setWebhooks(currentWebhooks)
}

func (p *Plex) removeWebhooks() {
	if p.Webhook == nil {
		return
	}
	for _, ip := range p.Webhook.ips {
		hookUrl := fmt.Sprintf("http://%s:%d/", ip.String(), p.Webhook.port)

		err := p.removeWebhook(hookUrl)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plex) removeWebhook(webhook string) error {
	currentWebhooks, err := p.getWebhooks()

	if err != nil {
		return err
	}

	for i, h := range currentWebhooks {
		if h == webhook {
			currentWebhooks = append(currentWebhooks[:i], currentWebhooks[i+1:]...)
			break
		}
	}

	return p.setWebhooks(currentWebhooks)
}

// SetWebhooks will set your webhooks to whatever you pass as an argument
// webhooks with a length of 0 will remove all webhooks
func (p *Plex) setWebhooks(webhooks []string) error {
	endpoint := "/api/v2/user/webhooks"

	body := url.Values{}

	if len(webhooks) == 0 {
		body.Add("urls[]", "")
	}

	for _, hook := range webhooks {
		body.Add("urls[]", hook)
	}

	err := PostHost(p, PlexURL, endpoint, []byte(body.Encode()))
	if err != nil {
		return err
	}
	return nil
}
