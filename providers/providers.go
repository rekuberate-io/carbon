package providers

type ProviderType string

const (
	WattTime        ProviderType = "watttime"
	ElectricityMaps ProviderType = "electricitymaps"
)

var (
	supportedProviders       = []ProviderType{WattTime, ElectricityMaps}
	supportedEmissionSignals = []EmissionsSignal{Average, Marginal}
)

func GetSupportedProviders() []ProviderType {
	return supportedProviders
}

type EmissionsSignal string

const (
	Average  EmissionsSignal = "average"
	Marginal EmissionsSignal = "marginal"
)

func GetSupportedEmissionSignals() []EmissionsSignal {
	return supportedEmissionSignals
}

type Provider interface {
	GetCarbonIntensity()
}

func NewProvider(providerType ProviderType, config ...string) (provider *Provider, err error) {
	switch providerType {
	case WattTime:
		provider, err = newWattTimeProvider()
		if err != nil {
			return nil, err
		}
	case ElectricityMaps:
		provider, err = newElectricityMapsProvider()
		if err != nil {
			return nil, err
		}
	}

	return provider, nil
}
