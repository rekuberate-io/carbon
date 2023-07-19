package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	wattTimeBaseUrl           string  = "https://api2.watttime.org/"
	wattTimeApiVersionUrlPath string  = "/v2"
	lbsTogramms               float64 = 453.59237
)

type WattTimeLiveResult struct {
	BalancingAuthority string    `json:"ba"`
	Frequency          string    `json:"freq"`
	Percent            string    `json:"percent"`
	MOER               string    `json:"moer"`
	PointTime          time.Time `json:"point_time"`
}

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
	//relativeLoginUrl := &url.URL{Path: "/v2/login"}
	//loginUrl := p.baseUrl.ResolveReference(relativeLoginUrl)

	loginUrl := ResolveAbsoluteUriReference(p.baseUrl, &url.URL{Path: wattTimeApiVersionUrlPath}, &url.URL{Path: "/login"})
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

func (p *WattTimeProvider) GetCurrent(ctx context.Context, zone *string) (string, error) {
	if zone == nil {
		return "", errors.New(fmt.Sprintf("zone (ba - balancing authority abbreviation) is required"))
	}

	requestUrl := ResolveAbsoluteUriReference(p.baseUrl, &url.URL{Path: wattTimeApiVersionUrlPath}, &url.URL{Path: "/index"})
	params := url.Values{}
	params.Add("ba", *zone)
	requestUrl.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token))
	response, err := p.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apierr, msg, err := p.unwrapHttpResponseErrorPayload(response)
		if err != nil {
			return "", errors.New(response.Status)
		}

		return "", errors.New(fmt.Sprintf("%s; %s: %s", response.Status, apierr, msg))
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result WattTimeLiveResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return "", err
	}

	moer, err := strconv.ParseFloat(result.MOER, 64)
	if err != nil {
		return "", nil
	}

	carbonIntensity := moer * lbsTogramms / 1000

	return fmt.Sprintf("%.2f", carbonIntensity), nil
}

func (p *WattTimeProvider) GetForecast(ctx context.Context, zone *string) (string, error) {
	return "", nil
}

func (p *WattTimeProvider) GetHistory(ctx context.Context, zone *string) (string, error) {
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
