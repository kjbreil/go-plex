package convert

import (
	"time"

	"github.com/kjbreil/go-plex/internal/plex/api"
	"github.com/kjbreil/go-plex/pkg/library"
)

// MetadataToShow converts API metadata to a library Show.
func MetadataToShow(m *api.MediaMetadata) *library.Show {
	for _, md := range m.MediaContainer.Metadata {
		if md.Type == "show" {
			return &library.Show{
				Title:          md.Title,
				Summary:        md.Summary,
				Year:           md.Year,
				ContentRating:  md.ContentRating,
				GUID:           md.GUID,
				TVDB:           md.AltGUIDs.TVDB(),
				Key:            "",
				RatingKey:      md.RatingKey,
				UserRating:     0,
				AudienceRating: 0,
				Watched:        false,
				LastViewedAt:   nil,
				AddedAt:        time.Time{},
				UpdatedAt:      time.Time{},
				Seasons:        nil,
				RefreshedAt:    time.Time{},
			}
		}
	}
	return nil
}

// MetadataToEpisode converts API metadata to a library Episode.
func MetadataToEpisode(m *api.MediaMetadata) *library.Episode {
	for _, md := range m.MediaContainer.Metadata {
		if md.Type == "episode" {
			return &library.Episode{
				Title:         md.Title,
				SeasonNumber:  int(md.Index),
				GUID:          md.GUID,
				TVDB:          0,
				ContentRating: md.ContentRating,
				Year:          md.Year,
				RatingKey:     md.RatingKey,
				Watched:       md.ViewCount > 0,
				LastViewedAt:  timeOrNil(md.LastViewedAt),
				AddedAt:       time.Unix(int64(md.AddedAt), 0),
				UpdatedAt:     time.Time{},
				RefreshedAt:   time.Time{},
				Duration:      0,
			}
		}
	}
	return nil
}

// UpdateShowFromMetadata updates a library Show with API metadata.
func UpdateShowFromMetadata(m *api.MediaMetadata, show *library.Show) {
	for _, md := range m.MediaContainer.Metadata {
		if md.GUID == show.GUID {
			show.Title = md.Title
			show.Year = md.Year
			show.Summary = md.Summary
			show.ContentRating = md.ContentRating
			show.RatingKey = md.RatingKey
			show.TVDB = md.AltGUIDs.TVDB()
			show.RefreshedAt = time.Now()
		}
	}
}

// UpdateEpisodeFromMetadata updates a library Episode with API metadata.
func UpdateEpisodeFromMetadata(m *api.MediaMetadata, ep *library.Episode) {
	for _, md := range m.MediaContainer.Metadata {
		if md.GUID == ep.GUID {
			ep.Title = md.Title
			ep.SeasonNumber = int(md.Index)
			ep.RatingKey = md.RatingKey
			ep.ContentRating = md.ContentRating
			ep.Year = md.Year
			ep.TVDB = md.AltGUIDs.TVDB()
			ep.Watched = md.ViewCount > 0
			ep.Duration = md.Duration
			ep.LastViewedAt = timeOrNil(md.LastViewedAt)
			ep.AddedAt = time.Unix(int64(md.AddedAt), 0)
			ep.RefreshedAt = time.Now()
		}
	}
}

// UpdateMovieFromMetadata updates a library Movie with API metadata.
func UpdateMovieFromMetadata(m *api.MediaMetadata, movie *library.Movie) {
	for _, md := range m.MediaContainer.Metadata {
		if md.GUID == movie.GUID {
			movie.Title = md.Title
			movie.Year = md.Year
			movie.RatingKey = md.RatingKey
			movie.ContentRating = md.ContentRating
			movie.Summary = md.Summary
			movie.TMDB = md.AltGUIDs.TMDB()
			movie.LastViewedAt = timeOrNil(md.LastViewedAt)
			movie.AddedAt = time.Unix(int64(md.AddedAt), 0)
			movie.RefreshedAt = time.Now()
		}
	}
}
