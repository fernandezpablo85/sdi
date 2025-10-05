package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fernandezpablo85/sdi/internal/api"
)

type HttpClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *HttpClient {
	return &HttpClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 1 * time.Second},
	}
}

func (c *HttpClient) GetAssetPrice(name string) (float64, error) {
	path := fmt.Sprintf("%s/v1/asset?name=%s", c.baseURL, name)
	res, err := c.httpClient.Get(path)
	if err != nil {
		return 0.0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0.0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body := api.AssetResponse{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return 0.0, err
	}
	return body.Price, nil
}

func (c *HttpClient) Ping() bool {
	path := fmt.Sprintf("%s/v1/healthz", c.baseURL)
	res, err := c.httpClient.Get(path)
	if err != nil {
		return false
	}
	if res.StatusCode != http.StatusOK {
		return false
	}
	defer res.Body.Close()
	return true
}
