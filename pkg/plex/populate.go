package plex

import (
	"log/slog"
	"sync"
	"time"

	"github.com/kjbreil/go-plex/pkg/library"
)

func (p *Plex) InitLibraries() error {
	var err error
	p.Libraries, err = p.GetLibraries()
	if err != nil {
		return err
	}

	p.mergeCache()

	return nil
}

func (p *Plex) PopulateLibraries() func() {
	done := make(chan struct{}, 1)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		start := time.Now()

		buf := make(chan struct{}, bufLen)
		wg := &sync.WaitGroup{}
		var err error
		for _, lib := range p.Libraries {
			if lib.Type == library.TypeShow {
				err = p.GetLibraryShows(lib, "")
				if err != nil {
					slog.Error("could not Get library shows", "library", lib.Title, "err", err.Error())
					continue
				}
			}
			if lib.Type == library.TypeMovie {
				err = p.GetLibraryMovies(lib, "")
				if err != nil {
					slog.Error("could not Get library movies", "library", lib.Title, "err", err.Error())
					continue
				}
			}
		}

		for _, lib := range p.Libraries {
			if lib.Type == library.TypeShow {
				for _, show := range lib.Shows {
					wg.Add(1)
					buf <- struct{}{}
					go func(show *library.Show) {
						defer func() {
							<-buf
							wg.Done()
						}()
						// exit out if the context is canceled
						if p.ctx.Err() != nil {
							return
						}
						err = p.GetShowEpisodes(show)
						if err != nil {
							slog.Error("could not Get show episodes", "show", show.Title, "err", err.Error())
							return
						}
					}(show)
				}
			}
		}

		wg.Wait()

		p.cleanupStaleLibraries(start)
		done <- struct{}{}
	}()

	return func() {
		<-done
		close(done)
	}
}

// cleanupStaleLibraries removes libraries, movies, shows, seasons, and episodes
// that were not refreshed since the given start time.
func (p *Plex) cleanupStaleLibraries(start time.Time) {
	libLength := len(p.Libraries)
	for i := 0; i < libLength; i++ {
		if p.Libraries[i].RefreshedAt.Before(start.Add(-1 * time.Minute)) {
			p.Libraries = append(p.Libraries[:i], p.Libraries[i+1:]...)
			i--
			libLength--
			continue
		}
		p.cleanupStaleMovies(p.Libraries[i], start)
		p.cleanupStaleShows(p.Libraries[i], start)
	}
}

func (p *Plex) cleanupStaleMovies(lib *library.Library, start time.Time) {
	moviesLength := len(lib.Movies)
	for j := 0; j < moviesLength; j++ {
		if lib.Movies[j].RefreshedAt.Before(start) {
			lib.Movies = append(lib.Movies[:j], lib.Movies[j+1:]...)
			j--
			moviesLength--
		}
	}
}

func (p *Plex) cleanupStaleShows(lib *library.Library, start time.Time) {
	showsLength := len(lib.Shows)
	for j := 0; j < showsLength; j++ {
		if lib.Shows[j].RefreshedAt.Before(start) {
			lib.Shows = append(lib.Shows[:j], lib.Shows[j+1:]...)
			j--
			showsLength--
			continue
		}
		for k, v := range lib.Shows[j].Seasons {
			if v.RefreshedAt.Before(start) {
				delete(lib.Shows[j].Seasons, k)
				continue
			}
			for l, e := range v.Episodes {
				if e.RefreshedAt.Before(start) {
					delete(lib.Shows[j].Seasons[k].Episodes, l)
				}
			}
		}
	}
}
