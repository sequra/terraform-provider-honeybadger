package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const HoneybadgerURL string = "https://app.honeybadger.io"

type HoneybadgerClient struct {
	HostURL    string
	HTTPClient *http.Client
	ApiToken   string
}

func NewClient(host *string, apiToken *string) *HoneybadgerClient {
	hbc := &HoneybadgerClient{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL:    HoneybadgerURL,
		ApiToken:   *apiToken,
	}

	if *host != "" {
		hbc.HostURL = *host
	}
	return hbc
}

func (hbc *HoneybadgerClient) DoRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(hbc.ApiToken, "")
	req.Header.Set("Content-Type", "application/json")

	res, err := hbc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) &&
		(res.StatusCode != http.StatusCreated) &&
		(res.StatusCode != http.StatusAccepted) &&
		(res.StatusCode != http.StatusNoContent) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
