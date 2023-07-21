package simulator

import (
	"context"
	"github.com/rekuberate-io/carbon/providers"
)

type Simulator struct {
}

func NewCarbonIntensityProviderSimulator() (*Simulator, error) {
	return &Simulator{}, nil
}

func (p *Simulator) GetCurrent(ctx context.Context, zone *string) (float64, error) {
	return providers.NoValue, nil
}

func (p *Simulator) GetForecast(ctx context.Context, zone *string) ([]providers.Forecast, error) {
	return nil, nil
}

func (p *Simulator) GetHistory(ctx context.Context, zone *string) (string, error) {
	return "", nil
}
