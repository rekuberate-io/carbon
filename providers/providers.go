package providers

import (
	"context"
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
	GetCurrent(ctx context.Context, zone *string) (string, error)
	GetForecast(ctx context.Context, zone *string) (string, error)
	GetHistory(ctx context.Context, zone *string) (string, error)
}

//
//type AbsoluteUriResolver interface {
//	resolveReference(baseUrl *url.URL, paths ...*url.URL) *url.URL
//}
