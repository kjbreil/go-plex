package plex

import (
	"github.com/kjbreil/go-plex/library"
	"log/slog"
	"sync"
	"time"
)

func (p *Plex) InitLibraries() error {
	var err error
	p.Libraries, err = p.GetLibraries()
	if err != nil {
		return err
	}

	p.mergeCache()

	// go p.PopulateLibraries()

	return nil
}

func (p *Plex) PopulateLibraries() func() {

	done := make(chan struct{}, 1)

	go func() {
		p.wg.Add(1)
		defer p.wg.Done()

		start := time.Now()

		buf := make(chan struct{}, bufLen)
		wg := &sync.WaitGroup{}
		var err error
		for _, lib := range p.Libraries {
			if lib.Type == library.TypeShow {
				err = p.GetLibraryShows(lib, "")
				if err != nil {
					slog.Error("could not Get library shows", "err", err.Error())
					continue
				}
				// wg.Add(len(lib.Shows))
			}
			if lib.Type == library.TypeMovie {
				err = p.GetLibraryMovies(lib, "")
				if err != nil {
					slog.Error("could not Get library movies", "err", err.Error())
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
							slog.Error("could not Get show episodes", "err", err.Error())
							return
						}
					}(show)
				}
			}
		}

		wg.Wait()

		// clean anything not updated between start and now
		libLength := len(p.Libraries)
		for i := 0; i < libLength; i++ {
			if p.Libraries[i].RefreshedAt.Before(start.Add(-1 * time.Minute)) {
				p.Libraries = append(p.Libraries[:i], p.Libraries[i+1:]...)
				i--
				libLength--
				continue
			}
			moviesLength := len(p.Libraries[i].Movies)
			for j := 0; j < moviesLength; j++ {
				if p.Libraries[i].Movies[j].RefreshedAt.Before(start) {
					p.Libraries[i].Movies = append(p.Libraries[i].Movies[:j], p.Libraries[i].Movies[j+1:]...)
					j--
					moviesLength--
					continue
				}
			}
			showsLength := len(p.Libraries[i].Shows)
			for j := 0; j < showsLength; j++ {
				if p.Libraries[i].Shows[j].RefreshedAt.Before(start) {
					p.Libraries[i].Shows = append(p.Libraries[i].Shows[:j], p.Libraries[i].Shows[j+1:]...)
					j--
					showsLength--
					continue
				}
				for k, v := range p.Libraries[i].Shows[j].Seasons {
					if v.RefreshedAt.Before(start) {
						delete(p.Libraries[i].Shows[j].Seasons, k)
						continue
					}
					for l, e := range v.Episodes {
						if e.RefreshedAt.Before(start) {
							delete(p.Libraries[i].Shows[j].Seasons[k].Episodes, l)
						}
					}
				}
			}
		}
		done <- struct{}{}
	}()

	return func() {
		<-done
		close(done)
	}
}
