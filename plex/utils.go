package plex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var bufLen = 10

type bufResp struct {
	resp *http.Response
	err  error
}
type bufReq struct {
	req  *http.Request
	resp chan bufResp
}

func get[T any](p *Plex, pa string, query url.Values) (T, error) {

	var rtn T

	u, _ := url.Parse(p.url.String())
	u.Path = path.Join(p.url.Path, pa)
	u.RawQuery = query.Encode()

	req, reqErr := http.NewRequest("GET", u.String(), nil)
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

	if err != nil {
		return rtn, err
	}

	if resp.StatusCode != http.StatusOK {
		return rtn, errors.New(resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(&rtn); err != nil {
		return rtn, err
	}

	return rtn, nil
}

func getSt[T any](p *Plex, query string) (T, error) {

	var rtn T

	u, _ := url.Parse(p.url.String())

	u.Path = path.Join(p.url.Path, query)

	req, reqErr := http.NewRequest("GET", u.String(), nil)
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

	if err != nil {
		return rtn, err
	}

	if resp.StatusCode != http.StatusOK {
		return rtn, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rtn, err
	}

	bs := string(body)
	fmt.Println(bs)

	return rtn, nil
}
