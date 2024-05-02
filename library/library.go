package library

import (
	"encoding/json"
	"fmt"
)

type Libraries []*Library

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
}

type LibraryType int

func (l *LibraryType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
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

func (l Libraries) Type(t LibraryType) Libraries {
	var nl Libraries
	for _, lib := range l {
		if lib.Type == t {
			nl = append(nl, lib)
		}
	}
	return nl
}
