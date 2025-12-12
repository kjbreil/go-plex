package plex

import (
	"encoding/json"
	"os"

	"github.com/kjbreil/go-plex/pkg/library"
)

func (p *Plex) mergeCache() {
	if p.cacheLibrary != "" {
		var cacheLibrary library.Libraries
		var err error
		var file []byte
		file, err = os.ReadFile(p.cacheLibrary)
		if err == nil {
			err = json.Unmarshal(file, &cacheLibrary)
			if err == nil {
			topLoop:
				for _, cl := range cacheLibrary {
					for i := range p.Libraries {
						if cl.Title == p.Libraries[i].Title {
							p.Libraries[i].Merge(cl)
							continue topLoop
						}
					}
					p.Libraries = append(p.Libraries, cl)
				}
			}
		}
	}
}

func (p *Plex) WriteCache() {
	if p.cacheLibrary != "" {
		b, err := json.Marshal(&p.Libraries)
		if err == nil {
			_ = os.WriteFile(p.cacheLibrary, b, 0644)
		}
	}
}
