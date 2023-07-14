package providers

type ProviderType string

const (
	WattTime        ProviderType = "watttime"
	ElectricityMaps ProviderType = "electricitymaps"
)

var (
	supportedProviders = []ProviderType{WattTime, ElectricityMaps}
)

func GetSupportedProviders() []ProviderType {
	return supportedProviders
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
