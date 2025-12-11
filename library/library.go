//go:generate stringer -type=LibraryType -trimprefix=LibraryType

package library

import (
	"encoding/json"
	"fmt"
	"time"
)

// Library shows plex library metadata
type Library struct {
	Location   []Location  `json:"Location"`
	Agent      string      `json:"agent"`
	AllowSync  bool        `json:"allowSync"`
	Art        string      `json:"art"`
	Composite  string      `json:"composite"`
	CreatedAt  int         `json:"createdAt"`
	Filter     bool        `json:"filters"`
	Key        string      `json:"key"`
	Language   string      `json:"language"`
	Refreshing bool        `json:"refreshing"`
	Scanner    string      `json:"scanner"`
	Thumb      string      `json:"thumb"`
	Title      string      `json:"title"`
	Type       LibraryType `json:"type"`
	UpdatedAt  int         `json:"updatedAt"`
	UUID       string      `json:"uuid"`
	Shows      Shows
	Movies     Movies

	RefreshedAt time.Time `json:"refreshedAt"`
}

type LibraryType int

func (l *LibraryType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		// try to parse as a int
		var i int
		if err := json.Unmarshal(b, &i); err != nil {
			return err
		}
		*l = LibraryType(i)
		return nil
	}
	switch s {
	case "show":
		*l = TypeShow
	case "movie":
		*l = TypeMovie
	default:
		return fmt.Errorf("unknown library type: %s", s)
	}
	return nil
}

const (
	TypeShow LibraryType = iota
	TypeMovie
)

// Location is the path of a plex server directory
type Location struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

func (l *Library) SetRefreshedAt() {
	l.RefreshedAt = time.Now()
	if l.Shows != nil {
		l.Shows.SetRefreshedAt()
	}
	if l.Movies != nil {
		l.Movies.SetRefreshedAt()
	}
}

func (l *Library) Merge(ml *Library) {
	if l.Movies != nil || ml.Movies != nil {
		if l.Movies != nil {
			l.Movies.Merge(&ml.Movies)
		} else {
			l.Movies = ml.Movies
		}
	}
	if l.Shows != nil || ml.Shows != nil {
		if l.Shows != nil {
			l.Shows.Merge(&ml.Shows)
		} else {
			l.Shows = ml.Shows
		}
	}
}
