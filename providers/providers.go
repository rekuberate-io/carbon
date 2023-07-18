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
	GetCurrent()
	GetForecast()
}

type BaseProvider struct {
	Signal EmissionsSignal
}

func NewProvider(providerType ProviderType, signal EmissionsSignal, config ...string) (provider Provider, err error) {
	switch providerType {
	case WattTime:
		provider, err = NewWattTimeProvider(signal)
		if err != nil {
			return nil, err
		}
	case ElectricityMaps:
		provider, err = NewElectricityMapsProvider(signal)
		if err != nil {
			return nil, err
		}
	}

	return provider, nil
}
