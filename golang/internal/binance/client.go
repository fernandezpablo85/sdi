package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type priceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", c.baseURL, symbol)
	res, err := c.httpClient.Get(url)
	if err != nil {
		return 0.0, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0.0, fmt.Errorf("unexpected status: %d", res.StatusCode)
	}
	var body priceResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return 0.0, err
	}
	return body.Price, nil
}
