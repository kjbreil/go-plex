package plex

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

const PlexURL = "https://plex.tv"

func get[T any](p *Plex, pa string, query url.Values) (T, error) {
	return getHost[T](p, p.url.String(), pa, query)
}

func getHost[T any](p *Plex, host string, pa string, query url.Values) (T, error) {
	var rtn T

	u, err := url.Parse(host)
	if err != nil {
		return rtn, err
	}
	u.Path = path.Join(u.Path, pa)
	u.RawQuery = query.Encode()

	req, reqErr := http.NewRequest(http.MethodGet, u.String(), nil)
	if reqErr != nil {
		return rtn, reqErr
	}
	req.Header = p.defaultHeaders

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return rtn, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return rtn, errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&rtn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return rtn, nil
		}
		return rtn, err
	}

	return rtn, nil
}

func postHost(p *Plex, host string, pa string, body []byte) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}
	u.Path = path.Join(p.url.Path, pa)

	req, reqErr := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if reqErr != nil {
		return reqErr
	}
	req.Header = p.defaultHeaders
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}

	return nil
}
