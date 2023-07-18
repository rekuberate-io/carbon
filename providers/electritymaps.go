package providers

type ElectricityMapsProvider struct {
}

func NewElectricityMapsProvider() (*ElectricityMapsProvider, error) {
	return &ElectricityMapsProvider{}, nil
}

func (p *ElectricityMapsProvider) GetCurrent(emissionsType EmissionsType) (string, error) {

	return "", nil
}

func (p *ElectricityMapsProvider) GetForecast(emissionsType EmissionsType) (string, error) {

	return "", nil
}
