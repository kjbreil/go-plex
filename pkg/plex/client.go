package plex

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/kjbreil/go-plex/internal/plex/api"
	"github.com/kjbreil/go-plex/internal/plex/convert"
	"github.com/kjbreil/go-plex/pkg/library"
)

const bufLen = 3

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

type Options func(*Plex)

func New(baseURL, token string, options ...Options) (*Plex, error) {
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

// GetLibraries of your Plex server.
func (p *Plex) GetLibraries() (library.Libraries, error) {
	resp, err := get[api.LibrarySections](p, "/library/sections", nil)
	if err != nil {
		return nil, err
	}

	resp.MediaContainer.Directory.SetRefreshedAt()

	return resp.MediaContainer.Directory, nil
}

// GetLibraryShows adds the shows to the Library.
func (p *Plex) GetLibraryShows(lib *library.Library, filter string) error {
	query := path.Join("/library/sections/", lib.Key, "all"+filter)
	resp, err := get[api.SearchResults](p, query, nil)
	lib.Shows.Merge(convert.SearchResultsToShows(&resp))
	var md api.MediaMetadata

	for _, show := range lib.Shows {
		md, err = p.GetMetadata(show.RatingKey)
		if err != nil {
			return err
		}

		convert.UpdateShowFromMetadata(&md, show)
	}

	return err
}

// GetLibraryMovies adds the movies to the Library.
func (p *Plex) GetLibraryMovies(lib *library.Library, filter string) error {
	query := path.Join("/library/sections/", lib.Key, "all"+filter)
	resp, err := get[api.SearchResults](p, query, nil)

	respMovies := convert.SearchResultsToMovies(&resp)
	lib.Movies.Merge(respMovies)
	var md api.MediaMetadata

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
				slog.Error("could not Get metadata", "ratingKey", ratingKey, "movie", movie.Title, "err", err.Error())
			} else {
				convert.UpdateMovieFromMetadata(&md, movie)
			}
		}(movie.RatingKey)
	}

	wg.Wait()

	return err
}

// GetSessions of devices currently consuming media.
func (p *Plex) GetSessions() (api.CurrentSessions, error) {
	urlPath := "/status/sessions"
	return get[api.CurrentSessions](p, urlPath, nil)
}

func (p *Plex) GetShowEpisodes(show *library.Show) error {
	if show == nil {
		return errors.New("no show provided")
	}
	query := path.Join("/library/metadata/", show.RatingKey, "children")

	resp, err := get[api.SearchResultsEpisode](p, query, nil)
	if err != nil {
		return err
	}
	show.Seasons.Merge(convert.EpisodeResultsToSeasons(&resp))

	for _, sea := range show.Seasons {
		query = path.Join("/library/metadata/", sea.RatingKey, "children")
		resp, err = get[api.SearchResultsEpisode](p, query, nil)
		if err != nil {
			p.logger.Error("could not get season metadata", "err", err.Error())
			continue
		}
		sea.Episodes.Merge(convert.EpisodeResultsToEpisodes(&resp))

		var md api.MediaMetadata

		for _, ep := range sea.Episodes {
			md, err = p.GetMetadata(ep.RatingKey)
			if err != nil {
				p.logger.Error("could not get episode metadata", "err", err.Error())
				continue
			}
			convert.UpdateEpisodeFromMetadata(&md, ep)
		}
	}

	return err
}

func (p *Plex) GetMetadata(ratingKey string) (api.MediaMetadata, error) {
	if ratingKey == "" {
		return api.MediaMetadata{}, errors.New("no ratingKey provided")
	}
	query := path.Join("/library/metadata/", ratingKey)

	resp, err := get[api.MediaMetadata](p, query, nil)

	return resp, err
}

type blank struct{}

func (p *Plex) Scrobble(key string) error {
	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := get[blank](p, "/:/scrobble", query)
	return err
}

func (p *Plex) UnScrobble(key string) error {
	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := get[blank](p, "/:/unscrobble", query)
	return err
}

func (p *Plex) ScanLibrary(lib *library.Library) error {
	urlPath := path.Join("/library/sections/", lib.Key, "refresh")

	for _, loc := range lib.Location {
		query := url.Values{}
		query.Add("path", loc.Path)
		_, err := get[blank](p, urlPath, query)
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
