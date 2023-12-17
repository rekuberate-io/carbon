package simulator

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/montanaflynn/stats"
	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
	"github.com/rekuberate-io/carbon/pkg/common"
	"time"
)

var (
	//go:embed latest.json
	latest string

	//go:embed forecast.json
	forecast string
)

type Simulator struct {
	bootstrap   bool
	replacement bool
}

func NewProvider(o carbonv1alpha1.Simulator) (*Simulator, error) {
	return &Simulator{bootstrap: *o.Spec.Bootstrap, replacement: *o.Spec.Replacement}, nil
}

func (p *Simulator) GetCurrent(ctx context.Context, zone string) (float64, error) {
	if p.bootstrap {
		sample, err := p.getSample(1, p.replacement)
		if err != nil {
			return common.NoValue, err
		}

		for _, point := range sample {
			return point, nil
		}
	}

	var result LiveResult
	err := json.Unmarshal([]byte(latest), &result)
	if err != nil {
		return common.NoValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *Simulator) GetForecast(ctx context.Context, zone string) (map[time.Time]float64, error) {
	return p.getSample(24, p.replacement)
}

func (p *Simulator) getSample(takenum int, replacement bool) (map[time.Time]float64, error) {
	var result ForecastResult
	err := json.Unmarshal([]byte(forecast), &result)
	if err != nil {
		return nil, err
	}

	forecasts := make(map[time.Time]float64)
	pointTime := time.Now()

	fd := []float64{}
	for _, f := range result.Forecast {
		fd = append(fd, float64(f.CarbonIntensity))
	}

	input := stats.Float64Data(fd)
	output, err := stats.Sample(input, takenum, replacement)
	if err != nil {
		return nil, err
	}

	for _, point := range output {
		forecasts[pointTime] = point
		pointTime = pointTime.Add(1 * time.Hour)
	}

	return forecasts, nil
}
