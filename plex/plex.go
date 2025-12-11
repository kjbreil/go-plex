package plex

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/go-plex/library"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"sync"
	"time"
)

type Plex struct {
	url            *url.URL
	token          string
	defaultHeaders http.Header
	httpClient     http.Client
	ctx            context.Context
	cancel         context.CancelFunc

	Libraries library.Libraries

	wg *sync.WaitGroup

	Websocket *NotificationEvents
	Webhook   *Webhook

	cacheLibrary string
	logger       *slog.Logger
}

type PlexOptions func(*Plex)

func New(baseURL, token string, options ...PlexOptions) (*Plex, error) {
	var p Plex

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.logger = slog.Default()
	p.wg = &sync.WaitGroup{}
	p.Websocket = NewNotificationEvents()

	if baseURL == "" && token == "" {
		return &p, errors.New("url or token is required")
	}

	p.httpClient = http.Client{
		Timeout: 3 * time.Second,
	}

	p.defaultHeaders = make(http.Header)
	p.defaultHeaders.Set("Content-Type", "application/json")
	p.defaultHeaders.Add("Accept", "application/json")
	p.defaultHeaders.Add("X-Plex-Platform", runtime.GOOS)
	p.defaultHeaders.Add("X-Plex-Platform-Version", "0.0.0")
	p.defaultHeaders.Add("X-Plex-Provides", "")
	p.defaultHeaders.Add("X-Plex-Client-Identifier", "go-plex-v0.0.1")
	p.defaultHeaders.Add("X-Plex-Product", "Go Plex")
	p.defaultHeaders.Add("X-Plex-Version", "0.0.1")
	p.defaultHeaders.Add("X-Plex-Device", runtime.GOOS+" "+runtime.GOARCH)
	p.defaultHeaders.Add("X-Plex-Token", token)

	for _, o := range options {
		o(&p)
	}

	var err error

	// has url and token
	if baseURL != "" && token != "" {
		p.url, err = url.ParseRequestURI(baseURL)

		p.token = token

		return &p, err
	}

	// just has token
	if baseURL == "" && token != "" {
		p.token = token

		return &p, nil
	}

	// just url
	p.url, err = url.ParseRequestURI(baseURL)

	return &p, err
}

// GetLibraries of your Plex server
func (p *Plex) GetLibraries() (library.Libraries, error) {
	resp, err := Get[LibrarySections](p, "/library/sections", nil)
	if err != nil {
		return nil, err
	}

	resp.MediaContainer.Directory.SetRefreshedAt()

	return resp.MediaContainer.Directory, nil
}

// GetLibraryShows adds the shows to the Library
func (p *Plex) GetLibraryShows(lib *library.Library, filter string) error {
	query := path.Join("/library/sections/", lib.Key, "all"+filter)
	resp, err := Get[SearchResults](p, query, nil)
	lib.Shows.Merge(resp.toShows())
	var md MediaMetadata

	for _, show := range lib.Shows {
		// p.BufGetMetadata(show.RatingKey, show)
		md, err = p.GetMetadata(show.RatingKey)
		if err != nil {
			return err
		}

		md.updateShow(show)
	}

	return err
}

// GetLibraryShows adds the shows to the Library
func (p *Plex) GetLibraryMovies(lib *library.Library, filter string) error {
	query := path.Join("/library/sections/", lib.Key, "all"+filter)
	resp, err := Get[SearchResults](p, query, nil)

	respMovies := resp.toMovies()
	lib.Movies.Merge(respMovies)
	var md MediaMetadata

	buf := make(chan struct{}, bufLen)
	wg := sync.WaitGroup{}

	wg.Add(len(lib.Movies))

	for _, movie := range lib.Movies {
		go func(ratingKey string) {

			buf <- struct{}{}
			defer func() {
				<-buf
				wg.Done()
			}()
			// exit out if the context is canceled
			if p.ctx.Err() != nil {
				return
			}
			md, err = p.GetMetadata(ratingKey)
			if err != nil {
				slog.Error("could not Get metadata", "err", err.Error())
			} else {
				md.updateMovie(movie)
			}

		}(movie.RatingKey)

	}

	wg.Wait()

	return err
}

// GetSessions of devices currently consuming media
func (p *Plex) GetSessions() (CurrentSessions, error) {
	path := "/status/sessions"
	return Get[CurrentSessions](p, path, nil)
}

func (p *Plex) GetShowEpisodes(show *library.Show) error {
	if show == nil {
		return fmt.Errorf("no show provided")
	}
	query := path.Join("/library/metadata/", show.RatingKey, "children")

	resp, err := Get[SearchResultsEpisode](p, query, nil)
	if err != nil {
		return err
	}
	show.Seasons.Merge(resp.toSeasons())

	for _, sea := range show.Seasons {

		query = path.Join("/library/metadata/", sea.RatingKey, "children")
		resp, err = Get[SearchResultsEpisode](p, query, nil)
		if err != nil {
			p.logger.Error("could not get season metadata", "err", err.Error())
			continue
		}
		sea.Episodes.Merge(resp.toEpisodes())

		var md MediaMetadata

		for _, ep := range sea.Episodes {
			md, err = p.GetMetadata(ep.RatingKey)
			if err != nil {
				p.logger.Error("could not get episode metadata", "err", err.Error())
				continue
			}
			md.updateEpisode(ep)
		}
	}

	return err
}

func (p *Plex) GetMetadata(ratingKey string) (MediaMetadata, error) {
	if ratingKey == "" {
		return MediaMetadata{}, errors.New("no ratingKey provided")
	}
	query := path.Join("/library/metadata/", ratingKey)

	resp, err := Get[MediaMetadata](p, query, nil)

	return resp, err
}

type blank struct{}

func (p *Plex) Scrobble(key string) error {
	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := Get[blank](p, "/:/", nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plex) UnScrobble(key string) error {
	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := Get[blank](p, "/:/", nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plex) ScanLibrary(lib *library.Library) error {
	urlPath := path.Join("/library/sections/", lib.Key, "refresh")

	for _, loc := range lib.Location {
		query := url.Values{}
		query.Add("path", loc.Path)
		_, err := Get[blank](p, urlPath, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plex) Close() {
	p.cancel()
	p.removeWebhooks()
	p.wg.Wait()
	p.WriteCache()
}
