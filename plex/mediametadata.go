package plex

import (
	"github.com/kjbreil/go-plex/library"
	"time"
)

type MediaMetadata struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
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
		}
	}
}

func (m *MediaMetadata) updateEpisode(ep *library.Episode) {
	for _, md := range m.MediaContainer.Metadata {
		if md.GUID == ep.GUID {
			ep.Title = md.Title
			ep.RatingKey = md.RatingKey
			ep.ContentRating = md.ContentRating
			ep.Year = md.Year
			ep.TVDB = md.AltGUIDs.TVDB()
			ep.Watched = md.ViewCount > 0
			ep.LastViewedAt = timeOrNil(md.LastViewedAt)
			ep.AddedAt = time.Unix(int64(md.AddedAt), 0)
		}
	}
}
