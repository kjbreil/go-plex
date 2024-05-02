package plex

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type boolOrInt struct {
	bool
}

func (b *boolOrInt) UnmarshalJSON(data []byte) error {
	var isInt int

	if err := json.Unmarshal(data, &isInt); err == nil {
		if isInt == 0 || isInt == 1 {

			if isInt != 0 && isInt != 1 {
				return fmt.Errorf("invalid boolOrInt: %d", isInt)
			}

			b.bool = isInt == 1

			return nil
		}
	}

	var isBool bool

	if err := json.Unmarshal(data, &isBool); err != nil {
		return err
	}

	b.bool = isBool

	return nil
}

type Ratings []Rating

type Rating struct {
	Image string  `json:"image"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func (r *Ratings) UnmarshalJSON(data []byte) error {
	float, err := strconv.ParseFloat(string(data), 64)
	if err == nil {
		*r = Ratings{
			Rating{
				Image: "",
				Type:  "",
				Value: float,
			},
		}
		return nil
	}
	var ratings []Rating
	if err = json.Unmarshal(data, &ratings); err == nil {
		*r = ratings
		return nil
	}

	return nil
}
