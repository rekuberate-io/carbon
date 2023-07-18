package providers

type WattTimeProvider struct {
	BaseProvider
	Username string
	ApiKey   string
}

func NewWattTimeProvider(signal EmissionsSignal, username string, apiKey string) (*WattTimeProvider, error) {
	return &WattTimeProvider{
		BaseProvider{Signal: signal},
		username,
		apiKey,
	}, nil
}

func (p *WattTimeProvider) GetCurrent() {

}

func (p *WattTimeProvider) GetForecast() {

}
