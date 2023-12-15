package electricitymaps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
	"github.com/rekuberate-io/carbon/pkg/common"
	"io"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

type ElectricityMapsProvider struct {
	subscription            SubscriptionType
	apiKey                  string
	zone                    string
	baseUrl                 *url.URL
	subscriptionRelativeUrl *url.URL
	client                  *http.Client
}

func NewProvider(ctx context.Context, k client.Client, o carbonv1alpha1.ElectricityMaps) (*ElectricityMapsProvider, error) {
	apiKeyRef := o.Spec.ApiKey
	objectKey := client.ObjectKey{
		Namespace: apiKeyRef.Namespace,
		Name:      apiKeyRef.Name,
	}
	secret := &corev1.Secret{}
	if err := k.Get(ctx, objectKey, secret); err != nil {
		return nil, err
	}

	apiKey := string(secret.Data["apiKey"])

	switch o.Spec.Subscription {
	case string(Commercial):
		return newElectricityMapsCommercialProvider(apiKey)
	case string(CommercialTrial):
		return newElectricityMapsCommercialTrialProvider(apiKey, o.Spec.CommercialTrialEndpoint)
	case string(FreeTier):
		return newElectricityMapsFreeTierProvider(apiKey)
	}

	return nil, nil
}

func newElectricityMapsCommercialProvider(apiKey string) (*ElectricityMapsProvider, error) {
	electricityMaps := &ElectricityMapsProvider{
		subscription: Commercial,
		apiKey:       apiKey,
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

func newElectricityMapsCommercialTrialProvider(apiKey string, commercialTrialEndpoint *string) (*ElectricityMapsProvider, error) {
	if commercialTrialEndpoint == nil {
		return nil, errors.New("no commercial trial id provided")
	}

	electricityMaps := &ElectricityMapsProvider{
		subscription:            CommercialTrial,
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

func newElectricityMapsFreeTierProvider(apiKey string) (*ElectricityMapsProvider, error) {
	electricityMaps := &ElectricityMapsProvider{
		subscription:            FreeTier,
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

func (p *ElectricityMapsProvider) GetCurrent(ctx context.Context) (float64, error) {
	requestUrl := common.ResolveAbsoluteUriReference(p.baseUrl, p.subscriptionRelativeUrl, &url.URL{Path: "/carbon-intensity/latest"})
	params := url.Values{}
	params.Add("zone", p.zone)
	requestUrl.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return common.NoValue, err
	}

	request.Header.Add("auth-token", p.apiKey)

	response, err := p.client.Do(request)
	if err != nil {
		return common.NoValue, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apierr, msg, err := p.unwrapHttpResponseErrorPayload(response)
		if err != nil {
			return common.NoValue, errors.New(response.Status)
		}

		return common.NoValue, errors.New(fmt.Sprintf("%s; %s: %s", response.Status, apierr, msg))
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return common.NoValue, err
	}

	var result LiveResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return common.NoValue, err
	}

	carbonIntensity := float64(result.CarbonIntensity)
	return carbonIntensity, nil
}

func (p *ElectricityMapsProvider) GetForecast(ctx context.Context) (map[time.Time]float64, error) {
	requestUrl := common.ResolveAbsoluteUriReference(
		p.baseUrl,
		p.subscriptionRelativeUrl,
		&url.URL{Path: "/carbon-intensity/forecast"},
	)
	params := url.Values{}
	params.Add("zone", p.zone)
	requestUrl.RawQuery = params.Encode()

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

	var result ForecastResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	forecasts := make(map[time.Time]float64)
	for _, f := range result.Forecast {
		//forecast := providers.Forecast{
		//	PointTime:       f.Datetime,
		//	CarbonIntensity: float64(f.CarbonIntensity),
		//}
		//
		//forecasts = append(forecasts, forecast)
		forecasts[f.Datetime] = float64(f.CarbonIntensity)
	}

	return forecasts, nil
}

func (p *ElectricityMapsProvider) Region() string {
	return p.zone
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
