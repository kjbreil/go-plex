package library

import "time"

type Shows []*Show

func (s *Shows) SetRefreshedAt() {
	for _, show := range *s {
		show.SetRefreshedAt()
	}
}
func (s *Shows) FindTitle(t string) *Show {
	for _, show := range *s {
		if show.Title == t {
			return show
		}
	}
	return nil
}

func (s *Shows) FindTvdbID(tvdb int) (*Show, *Season, *Episode) {
	for _, show := range *s {
		if show.TVDB == tvdb {
			return show, nil, nil
		}
		for _, season := range show.Seasons {
			for _, episode := range season.Episodes {
				if episode.TVDB == tvdb {
					return show, season, episode
				}
			}
		}
	}
	return nil, nil, nil
}

func (s *Shows) Merge(mergeShows *Shows) {
	if *s == nil {
		*s = *mergeShows
		return
	}
	for _, show := range *mergeShows {
		ss := s.FindTitle(show.Title)
		if ss != nil {
			ss.Merge(show)
		} else {
			*s = append(*s, show)
		}
	}
}

type Show struct {
	Title          string     `json:"title"`
	Summary        string     `json:"summary"`
	Year           int        `json:"year"`
	ContentRating  string     `json:"contentRating"`
	GUID           string     `json:"guid"`
	TVDB           int        `json:"tvdb"`
	Key            string     `json:"key"`
	RatingKey      string     `json:"ratingKey"`
	UserRating     float64    `json:"userRating"`
	AudienceRating float64    `json:"audienceRating"`
	Watched        bool       `json:"watched"`
	LastViewedAt   *time.Time `json:"lastViewedAt"`
	AddedAt        time.Time  `json:"addedAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Seasons        Seasons    `json:"seasons"`
	RefreshedAt    time.Time  `json:"refreshedAt"`
}

func (s *Show) SetRefreshedAt() {
	s.RefreshedAt = time.Now()
	s.Seasons.SetRefreshedAt()
}

func (s *Show) Merge(mergeShow *Show) {
	if s.RefreshedAt.Before(mergeShow.RefreshedAt) {
		mergeShow.Seasons.Merge(&s.Seasons)
		*s = *mergeShow
	}
	s.Seasons.Merge(&mergeShow.Seasons)
}
