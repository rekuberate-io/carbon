package watttime

import "time"

type LiveResult struct {
	BalancingAuthority string    `json:"ba"`
	Frequency          string    `json:"freq"`
	Percent            string    `json:"percent"`
	MOER               string    `json:"moer"`
	PointTime          time.Time `json:"point_time"`
}

type ForecastResult struct {
	GeneratedAt time.Time `json:"generated_at"`
	Forecast    []struct {
		PointTime time.Time `json:"point_time"`
		Value     float64   `json:"value"`
		Version   string    `json:"version"`
		Ba        string    `json:"ba"`
	} `json:"forecast"`
}
