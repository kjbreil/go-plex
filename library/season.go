package library

import "time"

type Seasons map[int]*Season

func (s *Seasons) SetRefreshedAt() {
	for _, season := range *s {
		season.SetRefreshedAt()
	}
}

func (s *Seasons) Merge(seasonsMerge *Seasons) {
	if *s == nil {
		*s = *seasonsMerge
		return
	}
	for k, season := range *seasonsMerge {
		ss := s.FindRatingKey(season.RatingKey)
		if ss != nil {
			ss.Merge(season)
		} else {
			(*s)[k] = season
		}
	}
}

func (s *Seasons) FindRatingKey(ratingKey string) *Season {
	for _, season := range *s {
		if season.RatingKey == ratingKey {
			return season
		}
	}
	return nil
}

type Season struct {
	Title     string `json:"title"`
	Number    int    `json:"number"`
	GUID      string `json:"guid"`
	RatingKey string `json:"ratingKey"`
	Episodes  Episodes

	RefreshedAt time.Time `json:"refreshedAt"`
}

func (s *Season) SetRefreshedAt() {
	s.RefreshedAt = time.Now()
	s.Episodes.SetRefreshedAt()
}

func (s *Season) Merge(seasonMerge *Season) {
	if s.RefreshedAt.Before(seasonMerge.RefreshedAt) {
		seasonMerge.Episodes.Merge(&s.Episodes)
		*s = *seasonMerge
	}
	s.Episodes.Merge(&seasonMerge.Episodes)
}
