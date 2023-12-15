package simulator

import "time"

type LiveResult struct {
	Zone               string    `json:"zone"`
	CarbonIntensity    int       `json:"carbonIntensity"`
	Datetime           time.Time `json:"datetime"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CreatedAt          time.Time `json:"createdAt"`
	EmissionFactorType string    `json:"emissionFactorType"`
	IsEstimated        bool      `json:"isEstimated"`
	EstimationMethod   string    `json:"estimationMethod"`
}

type ForecastResult struct {
	Zone     string `json:"zone"`
	Forecast []struct {
		CarbonIntensity int       `json:"carbonIntensity"`
		Datetime        time.Time `json:"datetime"`
	} `json:"forecast"`
	UpdatedAt time.Time `json:"updatedAt"`
}
