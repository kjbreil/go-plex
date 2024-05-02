package plex

import (
	"github.com/kjbreil/go-plex/library"
	"time"
)

// SearchResults ...
type SearchResults struct {
	MediaContainer SearchMediaContainer `json:"MediaContainer"`
}

func (s *SearchResults) toShows() library.Shows {
	shows := make(library.Shows, len(s.MediaContainer.Metadata))
	for i, m := range s.MediaContainer.Metadata {
		shows[i] = &library.Show{
			Title:          m.Title,
			Year:           m.Year,
			Summary:        m.Summary,
			ContentRating:  m.ContentRating,
			GUID:           m.GUID,
			TVDB:           m.AltGUIDs.TVDB(),
			Key:            m.Key,
			RatingKey:      m.RatingKey,
			UserRating:     m.UserRating,
			AudienceRating: m.AudienceRating,
			Watched:        m.ViewCount > 0,
			LastViewedAt:   timeOrNil(m.LastViewedAt),
			AddedAt:        time.Unix(int64(m.AddedAt), 0),
			UpdatedAt:      time.Unix(int64(m.UpdatedAt), 0),
		}
	}
	return shows
}

func timeOrNil(i int) *time.Time {
	if i == 0 {
		return nil
	}
	t := time.Unix(int64(i), 0)
	return &t
}
