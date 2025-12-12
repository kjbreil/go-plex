package plex

import (
	"encoding/json"
	"os"

	"github.com/kjbreil/go-plex/pkg/library"
)

func (p *Plex) mergeCache() {
	if p.cacheLibrary == "" {
		return
	}

	file, err := os.ReadFile(p.cacheLibrary)
	if err != nil {
		return
	}

	var cacheLibrary library.Libraries
	if err = json.Unmarshal(file, &cacheLibrary); err != nil {
		return
	}

	for _, cl := range cacheLibrary {
		merged := false
		for i := range p.Libraries {
			if cl.Title == p.Libraries[i].Title {
				p.Libraries[i].Merge(cl)
				merged = true
				break
			}
		}
		if !merged {
			p.Libraries = append(p.Libraries, cl)
		}
	}
}

func (p *Plex) WriteCache() {
	if p.cacheLibrary != "" {
		b, err := json.Marshal(&p.Libraries)
		if err == nil {
			_ = os.WriteFile(p.cacheLibrary, b, 0600)
		}
	}
}
