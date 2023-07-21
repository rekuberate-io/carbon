package providers

import (
	"context"
	"time"
)

type ProviderType string

const (
	WattTime             ProviderType = "watttime"
	ElectricityMaps      ProviderType = "electricitymaps"
	CarbonIntensityOrgUK ProviderType = "carbonintensity_org_uk"
)

var (
	supportedProviders     = []ProviderType{WattTime, ElectricityMaps, CarbonIntensityOrgUK}
	supportedEmissionTypes = []EmissionsType{Average, Marginal}
)

func GetSupportedProviders() []ProviderType {
	return supportedProviders
}

type EmissionsType string

const (
	Average  EmissionsType = "average"
	Marginal EmissionsType = "marginal"
)

func GetSupportedEmissionsTypes() []EmissionsType {
	return supportedEmissionTypes
}

type Provider interface {
	GetCurrent(ctx context.Context, zone *string) (float64, error)
	GetForecast(ctx context.Context, zone *string) ([]Forecast, error)
	GetHistory(ctx context.Context, zone *string) (string, error)
}

type Forecast struct {
	PointTime       time.Time `json:"pointTime"`
	CarbonIntensity float64   `json:"carbonIntensity"`
}
