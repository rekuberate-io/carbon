package providers

type ElectricityMapsProvider struct {
	BaseProvider
}

func NewElectricityMapsProvider(signal EmissionsSignal) (*ElectricityMapsProvider, error) {
	return &ElectricityMapsProvider{BaseProvider{Signal: signal}}, nil
}

func (p *ElectricityMapsProvider) GetCurrent() {

}

func (p *ElectricityMapsProvider) GetForecast() {

}
