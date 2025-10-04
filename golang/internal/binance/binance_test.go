package binance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPrice(t *testing.T) {
	asset := "BTCUSDT"
	expectedPrice := 43250.12
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"symbol":"%s","price":%f}`, asset, expectedPrice)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	price, err := client.GetPrice(asset)
	if err != nil {
		t.Fatalf("error while calling server %v", err)
	}
	if price != expectedPrice {
		t.Errorf("expected %f but was %f", expectedPrice, price)
	}
}
