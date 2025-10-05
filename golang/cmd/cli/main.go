package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/fernandezpablo85/sdi/internal/client"
	"github.com/fernandezpablo85/sdi/internal/env"
)

func main() {
	port := env.GetIntOrElse("PORT", 8080)
	baseUrl := fmt.Sprintf("http://localhost:%d", port)
	cli := client.NewClient(baseUrl)

	pong := cli.Ping()
	if !pong {
		log.Fatalf("could not ping api at %s", baseUrl)
	}

	slog.Info("ping success")

	assetName := "BTCUSDC"
	slog.Info("fetching asset price", "asset", assetName)
	res, err := cli.GetAssetPrice(assetName)
	if err != nil {
		log.Fatalf("error while fetching price: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("http error: %d", res.StatusCode)
	}
	slog.Info("asset price found", "asset", assetName, "price", res.Data.Price)
}
