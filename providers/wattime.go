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
	wattTimeBaseUrl string = "https://api2.watttime.org/"
)

type WattTimeProvider struct {
	baseUrl  *url.URL
	username string
	password string
	token    string
	client   *http.Client
}

func NewWattTimeProvider(ctx context.Context, username string, password string) (*WattTimeProvider, error) {
	watttime := &WattTimeProvider{client: &http.Client{
		Timeout: 10 * time.Second,
	}}

	baseUrl, err := url.Parse(wattTimeBaseUrl)
	if err != nil {
		return nil, err
	}

	watttime.baseUrl = baseUrl
	watttime.username = username
	watttime.password = password

	err = watttime.login(ctx)
	if err != nil {
		return nil, err
	}

	return watttime, nil
}

func (p *WattTimeProvider) login(ctx context.Context) error {
	relativeLoginUrl := &url.URL{Path: "/v2/login"}
	loginUrl := p.baseUrl.ResolveReference(relativeLoginUrl)
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

func (p *WattTimeProvider) GetCurrent(emissionsType EmissionsType) (string, error) {
	if emissionsType == Average {
		return "", errors.New(fmt.Sprintf("'%s' operating emissions are not supported", emissionsType))
	}

	return "", nil
}

func (p *WattTimeProvider) GetForecast(emissionsType EmissionsType) (string, error) {
	if emissionsType == Average {
		return "", errors.New(fmt.Sprintf("'%s' operating emissions are not supported", emissionsType))
	}

	return "", nil
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

//func (p *WattTimeProvider) GetMaps(ctx context.Context) (string, error) {
//	loginUrl := fmt.Sprintf("%s/%s", wattTimeBaseUrl, "maps")
//	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loginUrl, nil)
//	if err != nil {
//		return "", err
//	}
//
//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
//	resp, err := p.client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		apierr, msg, err := p.unwrapHttpResponseErrorPayload(resp)
//		if err != nil {
//			return "", errors.New(resp.Status)
//		}
//
//		return "", errors.New(fmt.Sprintf("%s, %s: %s", resp.Status, apierr, msg))
//	}
//
//	return "maps ok", nil
//}
