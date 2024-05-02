package plex

import (
	"errors"
	"fmt"
	"github.com/kjbreil/go-plex/library"
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
}

func New(baseURL, token string) (*Plex, error) {
	var p Plex

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

func (p *Plex) PopulateLibraries() (library.Libraries, error) {
	libraries, err := p.GetLibraries()
	if err != nil {
		return nil, err
	}

	buf := make(chan struct{}, 3)
	wg := sync.WaitGroup{}

	for _, lib := range libraries {
		if lib.Type == library.TypeShow {
			err = p.GetLibraryShows(lib, "")
			if err != nil {
				return nil, err
			}
			wg.Add(len(lib.Shows))
		}
	}

	for _, lib := range libraries {
		if lib.Type == library.TypeShow {
			for _, show := range lib.Shows {
				go func(show *library.Show) {
					buf <- struct{}{}
					defer func() {
						<-buf
						wg.Done()
					}()
					err = p.GetShowEpisodes(show)
					if err != nil {
						fmt.Println(err)
					}
				}(show)
				// err = p.GetShowEpisodes(show)
				// if err != nil {
				// 	return nil, err
				// }
				// wg.Done()
			}
		}
	}

	wg.Wait()

	return libraries, nil
}

// GetLibraries of your Plex server
func (p *Plex) GetLibraries() (library.Libraries, error) {
	resp, err := get[LibrarySections](p, "/library/sections", nil)
	if err != nil {
		return nil, err
	}

	return resp.MediaContainer.Directory, nil
}

// GetLibraryShows adds the shows to the Library
func (p *Plex) GetLibraryShows(lib *library.Library, filter string) error {
	query := path.Join("/library/sections/", lib.Key, "all"+filter)
	resp, err := get[SearchResults](p, query, nil)
	lib.Shows = resp.toShows()
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

func (p *Plex) GetShowEpisodes(show *library.Show) error {
	// start := time.Now()
	// fmt.Printf("%s started @ %s\n", show.Title, start.Format(time.TimeOnly))
	if show == nil {
		return fmt.Errorf("no show provided")
	}
	query := path.Join("/library/metadata/", show.RatingKey, "children")

	resp, err := get[SearchResultsEpisode](p, query, nil)
	if err != nil {
		return err
	}
	show.Seasons = resp.toSeasons()

	for _, sea := range show.Seasons {
		query = path.Join("/library/metadata/", sea.RatingKey, "children")
		resp, err = get[SearchResultsEpisode](p, query, nil)
		if err != nil {
			return err
		}
		sea.Episodes = resp.toEpisodes()

		var md MediaMetadata

		for _, ep := range sea.Episodes {
			// p.BufGetMetadata(ep.RatingKey, ep)

			md, err = p.GetMetadata(ep.RatingKey)
			if err != nil {
				return err
			}
			md.updateEpisode(ep)
		}
	}

	// fmt.Printf("%s finished @ %s took: %s\n", show.Title, time.Now().Format(time.TimeOnly), time.Now().Sub(start))
	return err
}

func (p *Plex) GetMetadata(ratingKey string) (MediaMetadata, error) {
	if ratingKey == "" {
		return MediaMetadata{}, errors.New("no ratingKey provided")
	}
	query := path.Join("/library/metadata/", ratingKey)

	resp, err := get[MediaMetadata](p, query, nil)

	return resp, err
}

type blank struct{}

func (p *Plex) Scrobble(key string) error {
	// query := fmt.Sprintf("/:/scrobble?key=%s&identifier=com.plexapp.plugins.library", key)

	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := get[blank](p, "/:/", nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plex) UnScrobble(key string) error {
	query := url.Values{}
	query.Add("key", key)
	query.Add("identifier", "com.plexapp.plugins.library")

	_, err := get[blank](p, "/:/", nil)
	if err != nil {
		return err
	}
	return nil
}
