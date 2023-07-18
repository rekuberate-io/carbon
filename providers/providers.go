package providers

type ProviderType string

const (
	WattTime        ProviderType = "watttime"
	ElectricityMaps ProviderType = "electricitymaps"
)

var (
	supportedProviders     = []ProviderType{WattTime, ElectricityMaps}
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
	GetCurrent(emissionsType EmissionsType) (string, error)
	GetForecast(emissionsType EmissionsType) (string, error)
}
