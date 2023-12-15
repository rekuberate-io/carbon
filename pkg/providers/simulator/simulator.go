package simulator

import (
	"context"
	_ "embed"
	"encoding/json"
	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
	"github.com/rekuberate-io/carbon/pkg/common"
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
	region    string
	max       float64
	min       float64
}

func NewProvider(o carbonv1alpha1.Simulator) (*Simulator, error) {
	randomize := *o.Spec.Randomize
	region := o.Spec.Region

	if randomize {
		var result ForecastResult
		err := json.Unmarshal([]byte(forecast), &result)
		if err != nil {
			return nil, err
		}

		mx, mn := getMaxMin(result)

		return &Simulator{
			region:    region,
			randomize: randomize,
			max:       float64(mx),
			min:       float64(mn),
		}, nil
	}

	return &Simulator{region: region, randomize: randomize}, nil
}

func (p *Simulator) GetCurrent(ctx context.Context) (float64, error) {
	if p.randomize {
		return rand.Float64() * (p.max - p.min), nil
	}

	var result electricitymaps.LiveResult
	err := json.Unmarshal([]byte(latest), &result)
	if err != nil {
		return common.NoValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *Simulator) GetForecast(ctx context.Context) (map[time.Time]float64, error) {
	var result ForecastResult
	err := json.Unmarshal([]byte(forecast), &result)
	if err != nil {
		return nil, err
	}

	forecasts := make(map[time.Time]float64)
	pointTime := time.Now()

	for range result.Forecast {
		pointTime = pointTime.Add(1 * time.Hour)

		//forecast := providers.Forecast{
		//	PointTime:       pointTime,
		//	CarbonIntensity: rand.Float64() * (p.max - p.min),
		//}
		//
		//forecasts = append(forecasts, forecast)

		forecasts[pointTime] = rand.Float64() * (p.max - p.min)
	}

	return forecasts, nil
}

func (p *Simulator) Region() string {
	return p.region
}

func getMaxMin(results ForecastResult) (int, int) {
	var mx int = results.Forecast[0].CarbonIntensity
	var mn int = results.Forecast[0].CarbonIntensity

	for _, value := range results.Forecast {
		if mx < value.CarbonIntensity {
			mx = value.CarbonIntensity
		}
		if mn > value.CarbonIntensity {
			mn = value.CarbonIntensity
		}
	}

	return mx, mn
}
