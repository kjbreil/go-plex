package library

import "time"

type Shows []*Show

func (s Shows) Title(t string) *Show {
	for _, show := range s {
		if show.Title == t {
			return show
		}
	}
	return nil
}

type Show struct {
	Title          string     `json:"title"`
	Summary        string     `json:"summary"`
	Year           int        `json:"year"`
	ContentRating  string     `json:"contentRating"`
	GUID           string     `json:"guid"`
	TVDB           string     `json:"tvdb"`
	Key            string     `json:"key"`
	RatingKey      string     `json:"ratingKey"`
	UserRating     float64    `json:"userRating"`
	AudienceRating float64    `json:"audienceRating"`
	Watched        bool       `json:"watched"`
	LastViewedAt   *time.Time `json:"lastViewedAt"`
	AddedAt        time.Time  `json:"addedAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Seasons        Seasons    `json:"seasons"`
}

type Seasons map[int]*Season

type Season struct {
	Title     string `json:"title"`
	GUID      string `json:"guid"`
	RatingKey string `json:"ratingKey"`
	Episodes  Episodes
}
type Episodes map[int]*Episode

type Episode struct {
	Title         string     `json:"title"`
	GUID          string     `json:"guid"`
	TVDB          string     `json:"tvdb"`
	ContentRating string     `json:"contentRating"`
	Year          int        `json:"year"`
	RatingKey     string     `json:"ratingKey"`
	Watched       bool       `json:"watched"`
	LastViewedAt  *time.Time `json:"lastViewedAt"`
	AddedAt       time.Time  `json:"addedAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
