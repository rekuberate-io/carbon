package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	electricityMapsBaseUrl      string = "https://api-access.electricitymaps.com/"
	electricityMapsFreeTierPath string = "/free-tier"
)

type SubscriptionType string

const (
	Commercial      SubscriptionType = "commercial"
	CommercialTrial SubscriptionType = "commercial_trial"
	FreeTier        SubscriptionType = "free_tier"
)

var (
	subscriptionTypes = []SubscriptionType{Commercial, CommercialTrial, FreeTier}
)

func GetElectricityMapSubscriptionModels() []SubscriptionType {
	return subscriptionTypes
}

type ElectricityMapLiveResult struct {
	Zone               string    `json:"zone"`
	CarbonIntensity    int       `json:"carbonIntensity"`
	Datetime           time.Time `json:"datetime"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CreatedAt          time.Time `json:"createdAt"`
	EmissionFactorType string    `json:"emissionFactorType"`
	IsEstimated        bool      `json:"isEstimated"`
	EstimationMethod   string    `json:"estimationMethod"`
}

type ElectricityMapForecastResult struct {
	Zone     string `json:"zone"`
	Forecast []struct {
		CarbonIntensity int       `json:"carbonIntensity"`
		Datetime        time.Time `json:"datetime"`
	} `json:"forecast"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ElectricityMapsProvider struct {
	subscription            SubscriptionType
	apiKey                  string
	baseUrl                 *url.URL
	subscriptionRelativeUrl *url.URL
	client                  *http.Client
}

func NewElectricityMapsProvider(apiKey string) (*ElectricityMapsProvider, error) {
	electricityMaps := &ElectricityMapsProvider{
		subscription: Commercial,
		//baseUrl:      &url.URL{Path: electricityMapsBaseUrl},
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	baseUrl, err := url.Parse(electricityMapsBaseUrl)
	if err != nil {
		return nil, err
	}

	electricityMaps.baseUrl = baseUrl

	return electricityMaps, nil
}

func NewElectricityMapsCommercialTrialProvider(apiKey string, commercialTrialEndpoint *string) (*ElectricityMapsProvider, error) {
	if commercialTrialEndpoint == nil {
		return nil, errors.New("no commercial trial id provided")
	}

	electricityMaps := &ElectricityMapsProvider{
		subscription: CommercialTrial,
		//baseUrl:                 &url.URL{Path: electricityMapsBaseUrl},
		subscriptionRelativeUrl: &url.URL{Path: *commercialTrialEndpoint},
		apiKey:                  apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	baseUrl, err := url.Parse(electricityMapsBaseUrl)
	if err != nil {
		return nil, err
	}

	electricityMaps.baseUrl = baseUrl

	return electricityMaps, nil
}

func NewElectricityMapsFreeTierProvider(apiKey string) (*ElectricityMapsProvider, error) {
	electricityMaps := &ElectricityMapsProvider{
		subscription: FreeTier,
		//baseUrl:                 &url.URL{Path: electricityMapsBaseUrl},
		subscriptionRelativeUrl: &url.URL{Path: electricityMapsFreeTierPath},
		apiKey:                  apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	baseUrl, err := url.Parse(electricityMapsBaseUrl)
	if err != nil {
		return nil, err
	}

	electricityMaps.baseUrl = baseUrl

	return electricityMaps, nil
}

func (p *ElectricityMapsProvider) GetCurrent(ctx context.Context, zone *string) (float64, error) {
	requestUrl := ResolveAbsoluteUriReference(p.baseUrl, p.subscriptionRelativeUrl, &url.URL{Path: "/carbon-intensity/latest"})

	if zone != nil {
		params := url.Values{}
		params.Add("zone", *zone)
		requestUrl.RawQuery = params.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return noValue, err
	}

	request.Header.Add("auth-token", p.apiKey)

	response, err := p.client.Do(request)
	if err != nil {
		return noValue, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apierr, msg, err := p.unwrapHttpResponseErrorPayload(response)
		if err != nil {
			return noValue, errors.New(response.Status)
		}

		return noValue, errors.New(fmt.Sprintf("%s; %s: %s", response.Status, apierr, msg))
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return noValue, err
	}

	var result ElectricityMapLiveResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return noValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *ElectricityMapsProvider) GetForecast(ctx context.Context, zone *string) ([]Forecast, error) {
	requestUrl := ResolveAbsoluteUriReference(
		p.baseUrl,
		p.subscriptionRelativeUrl,
		&url.URL{Path: "/carbon-intensity/forecast"},
	)

	if zone != nil {
		params := url.Values{}
		params.Add("zone", *zone)
		requestUrl.RawQuery = params.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("auth-token", p.apiKey)

	response, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apierr, msg, err := p.unwrapHttpResponseErrorPayload(response)
		if err != nil {
			return nil, errors.New(response.Status)
		}

		return nil, errors.New(fmt.Sprintf("%s; %s: %s", response.Status, apierr, msg))
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result ElectricityMapForecastResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	forecasts := make([]Forecast, 0)
	for _, f := range result.Forecast {
		forecast := Forecast{
			PointTime:       f.Datetime,
			CarbonIntensity: float64(f.CarbonIntensity),
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}

func (p *ElectricityMapsProvider) GetHistory(ctx context.Context, zone *string) (string, error) {
	return "", nil
}

func (p *ElectricityMapsProvider) unwrapHttpResponseErrorPayload(response *http.Response) (apiError string, message string, err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return response.Status, "", err
	}

	var errorPayload map[string]string
	err = json.Unmarshal(bytes, &errorPayload)
	if err != nil {
		return response.Status, "", err
	}

	return response.Status, errorPayload["error"], nil
}
