package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAssetPrice(t *testing.T) {
	tests := []struct {
		name           string
		assetName      string
		responseStatus int
		responseBody   string
		expectedPrice  float64
	}{
		{
			name:           "returns OK for existing asset",
			assetName:      "BTCUSDT",
			responseStatus: 200,
			responseBody:   `{"name": "BTCUSDT", "price": 4321.56}`,
			expectedPrice:  4321.56,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			res, err := client.GetAssetPrice(tt.assetName)
			if err != nil {
				t.Fatalf("error while calling server %v", err)
			}
			if res.StatusCode != tt.responseStatus {
				t.Fatalf("incorrect status, expected: %d got %d", tt.responseStatus, res.StatusCode)
			}
			if res.StatusCode == http.StatusOK && res.Data.Price != tt.expectedPrice {
				t.Errorf("expected %f but was %f", tt.expectedPrice, res.Data.Price)
			}
		})
	}
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	client := NewClient(server.URL)
	res := client.Ping()
	if !res {
		t.Errorf("ping returned false")
	}
}
