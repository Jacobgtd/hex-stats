package noaa

import (
	"encoding/json"
	"strconv"
	"time"
)

type StationsResponse struct {
	Stations []Station `json:"stations"`
}

type Station struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	State    string  `json:"state"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Timezone string  `json:"timezone"`
}

const noaaTimeLayout = "2006-01-02 15:04"

type TidePredictionsResponse struct {
	Predictions []TidePrediction `json:"predictions"`
}

type TideType string

const (
	HighTide TideType = "H"
	LowTide  TideType = "L"
)

type TidePrediction struct {
	Time  time.Time `json:"t"`
	Value float64   `json:"v"`
	Type  TideType  `json:"type"`
}

func (p *TidePrediction) UnmarshalJSON(data []byte) error {
	type rawPrediction struct {
		Time  string   `json:"t"`
		Value string   `json:"v"`
		Type  TideType `json:"type"`
	}

	var raw rawPrediction
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, err := time.Parse(noaaTimeLayout, raw.Time)
	if err != nil {
		return err
	}

	v, err := strconv.ParseFloat(raw.Value, 64)
	if err != nil {
		return err
	}

	p.Time = t
	p.Value = v
	p.Type = raw.Type

	return nil
}
