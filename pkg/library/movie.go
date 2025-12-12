package library

import "time"

type Movies []*Movie

func (m *Movies) SetRefreshedAt() {
	for _, movie := range *m {
		movie.SetRefreshedAt()
	}
}

func (m *Movies) Merge(mergeMovies *Movies) {
	if *m == nil {
		*m = *mergeMovies
		return
	}
	for _, movie := range *mergeMovies {
		mm := m.FindRatingKey(movie.RatingKey)
		if mm != nil {
			mm.Merge(movie)
		} else {
			*m = append(*m, movie)
		}
	}
}

func (m Movies) FindRatingKey(ratingKey string) *Movie {
	for _, movie := range m {
		if movie.RatingKey == ratingKey {
			return movie
		}
	}
	return nil
}

func (m *Movies) FindTMDB(id int) *Movie {
	for _, movie := range *m {
		if movie.TMDB == id {
			return movie
		}
	}
	return nil
}

type Movie struct {
	Title          string     `json:"title"`
	Summary        string     `json:"summary"`
	Year           int        `json:"year"`
	ContentRating  string     `json:"contentRating"`
	GUID           string     `json:"guid"`
	TMDB           int        `json:"tmdb"`
	Key            string     `json:"key"`
	RatingKey      string     `json:"ratingKey"`
	UserRating     float64    `json:"userRating"`
	AudienceRating float64    `json:"audienceRating"`
	Watched        bool       `json:"watched"`
	LastViewedAt   *time.Time `json:"lastViewedAt"`
	AddedAt        time.Time  `json:"addedAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`

	RefreshedAt time.Time `json:"refreshedAt"`
}

func (m *Movie) SetRefreshedAt() {
	m.RefreshedAt = time.Now()
}

func (m *Movie) Merge(mm *Movie) {
	if m.RefreshedAt.Before(mm.RefreshedAt) {
		*m = *mm
	}
}
