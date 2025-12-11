package plex

import (
	"github.com/kjbreil/go-plex/library"
	"time"
)

// SearchResultsEpisode contains metadata about an episode
type SearchResultsEpisode struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

func (s *SearchResultsEpisode) toSeasons() *library.Seasons {
	seasons := make(library.Seasons, len(s.MediaContainer.Metadata))
	for _, m := range s.MediaContainer.Metadata {
		seasons[int(m.Index)] = &library.Season{
			Title:       m.Title,
			Number:      int(m.Index),
			GUID:        m.GUID,
			RatingKey:   m.RatingKey,
			RefreshedAt: time.Now(),
		}
	}

	return &seasons
}

func (s *SearchResultsEpisode) toEpisodes() *library.Episodes {
	episodes := make(library.Episodes, len(s.MediaContainer.Metadata))
	for _, m := range s.MediaContainer.Metadata {
		episodes[int(m.Index)] = &library.Episode{
			Title:         m.Title,
			SeasonNumber:  int(m.Index),
			GUID:          m.GUID,
			RatingKey:     m.RatingKey,
			ContentRating: m.ContentRating,
			Year:          m.Year,
			Watched:       m.ViewCount > 0,
			Duration:      m.Duration,
			LastViewedAt:  timeOrNil(m.LastViewedAt),
			AddedAt:       time.Unix(int64(m.AddedAt), 0),
			UpdatedAt:     time.Unix(int64(m.UpdatedAt), 0),
			RefreshedAt:   time.Now(),
		}
	}

	return &episodes
}
