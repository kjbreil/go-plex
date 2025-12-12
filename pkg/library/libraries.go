package library

type Libraries []*Library

func (l Libraries) SetRefreshedAt() {
	for _, lib := range l {
		lib.SetRefreshedAt()
	}
}

func (l Libraries) Type(t LibraryType) Libraries {
	var nl Libraries
	for _, lib := range l {
		if lib.Type == t {
			nl = append(nl, lib)
		}
	}
	return nl
}

// FindEpisode gets the show season and epsidoe for a RatingKey.
func (l Libraries) FindEpisode(ratingKey string) (*Show, *Season, *Episode) {
	for _, lib := range l {
		for _, show := range lib.Shows {
			for _, season := range show.Seasons {
				for _, episode := range season.Episodes {
					if episode.RatingKey == ratingKey {
						return show, season, episode
					}
				}
			}
		}
	}
	return nil, nil, nil
}

// FindSeason gets the show season for a RatingKey for the season itself or RatingKey of an episode in the season.
func (l Libraries) FindSeason(ratingKey string) (*Show, *Season) {
	for _, lib := range l {
		for _, show := range lib.Shows {
			for _, season := range show.Seasons {
				if season.RatingKey == ratingKey {
					return show, season
				}
				for _, episode := range season.Episodes {
					if episode.RatingKey == ratingKey {
						return show, season
					}
				}
			}
		}
	}
	return nil, nil
}

// FindShow gets the show for a RatingKey for the show itself or RatingKey of a season contained in the show or an
// episode contained in the show.
func (l Libraries) FindShow(ratingKey string) *Show {
	for _, lib := range l {
		for _, show := range lib.Shows {
			if show.RatingKey == ratingKey {
				return show
			}
			for _, season := range show.Seasons {
				if season.RatingKey == ratingKey {
					return show
				}
				for _, episode := range season.Episodes {
					if episode.RatingKey == ratingKey {
						return show
					}
				}
			}
		}
	}
	return nil
}
