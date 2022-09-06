package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const HoneyBadgerURL string = "https://app.honeybadger.io"

type HoneyBadgerClient struct {
	HostURL    string
	HTTPClient *http.Client
	ApiToken   string
	TeamID     int
}

func NewClient(host *string, apiToken *string, teamID *int) *HoneyBadgerClient {
	hbc := &HoneyBadgerClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HoneyBadgerURL,
		ApiToken:   *apiToken,
		TeamID:     *teamID,
	}

	if *host != "" {
		hbc.HostURL = *host
	}
	return hbc
}

func (hbc *HoneyBadgerClient) doRequest(req *http.Request) ([]byte, error) {
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
		(res.StatusCode != http.StatusAccepted) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
