package library

import "time"

type Episodes map[int]*Episode

func (e *Episodes) SetRefreshedAt() {
	for _, episode := range *e {
		episode.SetRefreshedAt()
	}
}

func (e *Episodes) Merge(episodesMerge *Episodes) {
	if *e == nil {
		*e = *episodesMerge
		return
	}
	for k, episode := range *episodesMerge {
		ee := e.FindRatingKey(episode.RatingKey)
		if ee != nil {
			episode.Merge(ee)
		} else {
			(*e)[k] = episode
		}
	}
}

func (e *Episodes) FindRatingKey(ratingKey string) *Episode {
	for _, episode := range *e {
		if episode.RatingKey == ratingKey {
			return episode
		}
	}
	return nil
}

type Episode struct {
	Title         string     `json:"title"`
	SeasonNumber  int        `json:"seasonNumber"`
	GUID          string     `json:"guid"`
	TVDB          int        `json:"tvdb"`
	ContentRating string     `json:"contentRating"`
	Year          int        `json:"year"`
	RatingKey     string     `json:"ratingKey"`
	Watched       bool       `json:"watched"`
	LastViewedAt  *time.Time `json:"lastViewedAt"`
	AddedAt       time.Time  `json:"addedAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	RefreshedAt   time.Time  `json:"refreshedAt"`
	Duration      int        `json:"duration"`
}

func (e *Episode) SetRefreshedAt() {
	e.RefreshedAt = time.Now()
}

func (e *Episode) Merge(ee *Episode) {
	if e.RefreshedAt.Before(ee.RefreshedAt) {
		*e = *ee
	}
}
