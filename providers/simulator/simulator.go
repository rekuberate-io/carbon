package simulator

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/rekuberate-io/carbon/providers"
)

var (
	//go:embed latest.json
	latest string

	//go:embed forecast.json
	forecast string
)

const Zone = "SIM-1"

type Simulator struct {
}

func NewCarbonIntensityProviderSimulator() (*Simulator, error) {
	return &Simulator{}, nil
}

func (p *Simulator) GetCurrent(ctx context.Context, zone *string) (float64, error) {
	var result providers.ElectricityMapLiveResult
	err := json.Unmarshal([]byte(latest), &result)
	if err != nil {
		return providers.NoValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *Simulator) GetForecast(ctx context.Context, zone *string) ([]providers.Forecast, error) {
	var result providers.ElectricityMapForecastResult
	err := json.Unmarshal([]byte(forecast), &result)
	if err != nil {
		return nil, err
	}

	forecasts := make([]providers.Forecast, 0)
	for _, f := range result.Forecast {
		forecast := providers.Forecast{
			PointTime:       f.Datetime,
			CarbonIntensity: float64(f.CarbonIntensity),
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}

func (p *Simulator) GetHistory(ctx context.Context, zone *string) (string, error) {
	return "", nil
}
