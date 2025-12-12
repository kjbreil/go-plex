package convert

import (
	"time"

	"github.com/kjbreil/go-plex/internal/plex/api"
	"github.com/kjbreil/go-plex/pkg/library"
)

// SearchResultsToShows converts API search results to library Shows.
func SearchResultsToShows(s *api.SearchResults) *library.Shows {
	shows := make(library.Shows, len(s.MediaContainer.Metadata))
	for i, m := range s.MediaContainer.Metadata {
		shows[i] = &library.Show{
			Title:          m.Title,
			Year:           m.Year,
			Summary:        m.Summary,
			ContentRating:  m.ContentRating,
			GUID:           m.GUID,
			Key:            m.Key,
			RatingKey:      m.RatingKey,
			UserRating:     m.UserRating,
			AudienceRating: m.AudienceRating,
			Watched:        m.ViewCount > 0,
			LastViewedAt:   timeOrNil(m.LastViewedAt),
			AddedAt:        time.Unix(int64(m.AddedAt), 0),
			UpdatedAt:      time.Unix(int64(m.UpdatedAt), 0),
			RefreshedAt:    time.Now(),
		}
	}
	return &shows
}

// SearchResultsToMovies converts API search results to library Movies.
func SearchResultsToMovies(s *api.SearchResults) *library.Movies {
	movies := make(library.Movies, len(s.MediaContainer.Metadata))
	for i, m := range s.MediaContainer.Metadata {
		movies[i] = &library.Movie{
			Title:          m.Title,
			Year:           m.Year,
			Summary:        m.Summary,
			ContentRating:  m.ContentRating,
			GUID:           m.GUID,
			Key:            m.Key,
			RatingKey:      m.RatingKey,
			UserRating:     m.UserRating,
			AudienceRating: m.AudienceRating,
			Watched:        m.ViewCount > 0,
			LastViewedAt:   timeOrNil(m.LastViewedAt),
			AddedAt:        time.Unix(int64(m.AddedAt), 0),
			UpdatedAt:      time.Unix(int64(m.UpdatedAt), 0),
			RefreshedAt:    time.Now(),
		}
	}
	return &movies
}

// EpisodeResultsToSeasons converts API episode results to library Seasons.
func EpisodeResultsToSeasons(s *api.SearchResultsEpisode) *library.Seasons {
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

// EpisodeResultsToEpisodes converts API episode results to library Episodes.
func EpisodeResultsToEpisodes(s *api.SearchResultsEpisode) *library.Episodes {
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

func timeOrNil(i int) *time.Time {
	if i == 0 {
		return nil
	}
	t := time.Unix(int64(i), 0)
	return &t
}
