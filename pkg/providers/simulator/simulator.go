package simulator

import (
	"context"
	_ "embed"
	"encoding/json"
	providers2 "github.com/rekuberate-io/carbon/pkg/providers"
	"github.com/rekuberate-io/carbon/pkg/providers/electricitymaps"
	"math/rand"
	"time"
)

var (
	//go:embed latest.json
	latest string

	//go:embed forecast.json
	forecast string
)

type Simulator struct {
	randomize bool
	zone      string
	max       float64
	min       float64
}

func NewCarbonIntensityProviderSimulator(zone string, randomize bool) (*Simulator, error) {
	if randomize {
		var result electricitymaps.ElectricityMapForecastResult
		err := json.Unmarshal([]byte(forecast), &result)
		if err != nil {
			return nil, err
		}

		max, min := getMaxMin(result)

		return &Simulator{
			zone:      zone,
			randomize: randomize,
			max:       float64(max),
			min:       float64(min),
		}, nil
	}

	return &Simulator{zone: zone, randomize: randomize}, nil
}

func (p *Simulator) GetCurrent(ctx context.Context, zone string) (float64, error) {
	if p.randomize {
		return rand.Float64() * (p.max - p.min), nil
	}

	var result electricitymaps.ElectricityMapLiveResult
	err := json.Unmarshal([]byte(latest), &result)
	if err != nil {
		return providers2.NoValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *Simulator) GetForecast(ctx context.Context, zone string) ([]providers2.Forecast, error) {
	var result electricitymaps.ElectricityMapForecastResult
	err := json.Unmarshal([]byte(forecast), &result)
	if err != nil {
		return nil, err
	}

	forecasts := make([]providers2.Forecast, 0)
	pointTime := time.Now()

	for range result.Forecast {
		pointTime = pointTime.Add(1 * time.Hour)

		forecast := providers2.Forecast{
			PointTime:       pointTime,
			CarbonIntensity: rand.Float64() * (p.max - p.min),
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}

func getMaxMin(results electricitymaps.ElectricityMapForecastResult) (int, int) {
	var max int = results.Forecast[0].CarbonIntensity
	var min int = results.Forecast[0].CarbonIntensity

	for _, value := range results.Forecast {
		if max < value.CarbonIntensity {
			max = value.CarbonIntensity
		}
		if min > value.CarbonIntensity {
			min = value.CarbonIntensity
		}
	}

	return max, min
}
