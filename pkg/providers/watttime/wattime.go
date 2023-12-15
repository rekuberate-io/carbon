package watttime

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
	"strconv"
	"time"
)

const (
	wattTimeBaseUrl           string  = "https://api2.watttime.org/"
	wattTimeApiVersionUrlPath string  = "/v2"
	lbsTogramms               float64 = 453.59237
)

type WattTimeProvider struct {
	baseUrl  *url.URL
	username string
	password string
	region   string
	token    string
	client   *http.Client
}

func NewProvider(ctx context.Context, k client.Client, o carbonv1alpha1.WattTime) (*WattTimeProvider, error) {
	watttime := &WattTimeProvider{client: &http.Client{
		Timeout: 10 * time.Second,
	}}

	baseUrl, err := url.Parse(wattTimeBaseUrl)
	if err != nil {
		return nil, err
	}

	passwordRef := o.Spec.Password
	objectKey := client.ObjectKey{
		Namespace: passwordRef.Namespace,
		Name:      passwordRef.Name,
	}
	secret := &corev1.Secret{}
	if err := k.Get(ctx, objectKey, secret); err != nil {
		return nil, err
	}

	watttime.baseUrl = baseUrl
	watttime.username = o.Spec.Username
	watttime.password = string(secret.Data["password"])

	err = watttime.login(ctx)
	if err != nil {
		return nil, err
	}

	return watttime, nil
}

func (p *WattTimeProvider) login(ctx context.Context) error {
	loginUrl := common.ResolveAbsoluteUriReference(p.baseUrl, &url.URL{Path: wattTimeApiVersionUrlPath}, &url.URL{Path: "/login"})
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, loginUrl.String(), nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(p.username, p.password)
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apierr, msg, err := p.unwrapHttpResponseErrorPayload(response)
		if err != nil {
			return errors.New(response.Status)
		}

		return errors.New(fmt.Sprintf("%s; %s: %s", response.Status, apierr, msg))
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var tokenAsJson map[string]string
	err = json.Unmarshal(bytes, &tokenAsJson)
	if err != nil {
		return err
	}

	p.token = tokenAsJson["token"]

	return nil
}

func (p *WattTimeProvider) GetCurrent(ctx context.Context) (float64, error) {
	requestUrl := common.ResolveAbsoluteUriReference(p.baseUrl, &url.URL{Path: wattTimeApiVersionUrlPath}, &url.URL{Path: "/index"})
	params := url.Values{}
	params.Add("ba", p.region)
	requestUrl.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return common.NoValue, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token))
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

	moer, err := strconv.ParseFloat(result.MOER, 64)
	if err != nil {
		return common.NoValue, nil
	}

	carbonIntensity := moer * lbsTogramms / 1000
	return carbonIntensity, nil
}

func (p *WattTimeProvider) GetForecast(ctx context.Context) (map[time.Time]float64, error) {
	requestUrl := common.ResolveAbsoluteUriReference(
		p.baseUrl,
		&url.URL{Path: wattTimeApiVersionUrlPath},
		&url.URL{Path: "/forecast"},
	)
	params := url.Values{}
	params.Add("ba", p.region)
	requestUrl.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token))
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
		forecasts[f.PointTime] = f.Value
	}

	return forecasts, nil
}

func (p *WattTimeProvider) Region() string {
	return p.region
}

func (p *WattTimeProvider) unwrapHttpResponseErrorPayload(response *http.Response) (apiError string, message string, err error) {
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var errorPayload map[string]string
	err = json.Unmarshal(bytes, &errorPayload)
	if err != nil {
		return "", "", err
	}

	return errorPayload["error"], errorPayload["message"], nil
}
