package providers

import (
	"context"
	"slices"
	"strings"
	"time"
)

type ProviderType string

const (
	WattTime        ProviderType = "watttime"
	ElectricityMaps ProviderType = "electricitymaps"
	Simulator       ProviderType = "simulator"
)

var (
	supportedProviders     = []ProviderType{WattTime, ElectricityMaps, Simulator}
	supportedEmissionTypes = []EmissionsType{Average, Marginal}
)

type Provider interface {
	GetCurrent(ctx context.Context, zone string) (float64, error)
	GetForecast(ctx context.Context, zone string) ([]Forecast, error)
}

type Forecast struct {
	PointTime       time.Time `json:"pointTime"`
	CarbonIntensity float64   `json:"carbonIntensity"`
}

func GetSupportedProviders() []ProviderType {
	return supportedProviders
}

func IsSupported(providerType string) bool {
	if !slices.Contains(supportedProviders, ProviderType(strings.ToLower(providerType))) {
		return false
	}

	return true
}
