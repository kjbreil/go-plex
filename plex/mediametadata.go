package plex

import (
	"github.com/kjbreil/go-plex/library"
	"time"
)

type MediaMetadata struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

func (m *MediaMetadata) Show() *library.Show {
	for _, md := range m.MediaContainer.Metadata {
		if md.Type == "show" {
			return &library.Show{
				Title:         md.Title,
				Year:          md.Year,
				Summary:       md.Summary,
				ContentRating: md.ContentRating,
				RatingKey:     md.RatingKey,
				TVDB:          md.AltGUIDs.TVDB(),
				GUID:          md.GUID,
			}
		}
	}
	return nil
}

func (m *MediaMetadata) Episode() *library.Episode {
	for _, md := range m.MediaContainer.Metadata {
		if md.Type == "episode" {
			return &library.Episode{
				Title:         md.Title,
				SeasonNumber:  int(md.Index),
				GUID:          md.GUID,
				RatingKey:     md.RatingKey,
				ContentRating: md.ContentRating,
				Year:          md.Year,
				Watched:       md.ViewCount > 0,
				LastViewedAt:  timeOrNil(md.LastViewedAt),
				AddedAt:       time.Unix(int64(md.AddedAt), 0),
			}
		}
	}
	return nil
}

func (m *MediaMetadata) updateShow(show *library.Show) {
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

func (m *MediaMetadata) updateEpisode(ep *library.Episode) {
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

func (m *MediaMetadata) updateMovie(movie *library.Movie) {
	for _, md := range m.MediaContainer.Metadata {
		if md.GUID == movie.GUID {
			movie.Title = md.Title
			movie.Year = md.Year
			movie.RatingKey = md.RatingKey
			movie.ContentRating = md.ContentRating
			movie.Summary = md.Summary
			movie.ContentRating = md.ContentRating
			movie.RatingKey = md.RatingKey
			movie.TMDB = md.AltGUIDs.TMDB()
			movie.LastViewedAt = timeOrNil(md.LastViewedAt)
			movie.AddedAt = time.Unix(int64(md.AddedAt), 0)
			movie.RefreshedAt = time.Now()
		}
	}
}
