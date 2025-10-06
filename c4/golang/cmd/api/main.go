package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fernandezpablo85/sdi/internal/api"
	"github.com/fernandezpablo85/sdi/internal/binance"
	"github.com/fernandezpablo85/sdi/internal/env"
)

const DEFAULT_PORT = 8080

type Pricer interface {
	GetPrice(assetName string) (float64, error)
}

func only(method string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		f(w, r)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func assetHandler(client Pricer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			msg := "missing asset name"
			slog.Debug(msg, "status", http.StatusBadRequest)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		slog.Info("fetching asset price", "asset", name)
		price, err := client.GetPrice(name)
		if errors.Is(err, api.ErrAssetNotFound) {
			msg := "asset unknown"
			slog.Error(msg, "asset", name, "error", err)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		if err != nil {
			msg := "error fetching asset price"
			slog.Error(msg, "asset", name, "error", err)
			http.Error(w, msg, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.AssetResponse{Name: name, Price: price})
	}
}

func setupRoutes(mux *http.ServeMux, pricer Pricer) {
	mux.HandleFunc("/v1/healthz", only(http.MethodGet, healthzHandler))
	mux.HandleFunc("/v1/asset", only(http.MethodGet, assetHandler(pricer)))
}

func main() {
	mux := http.NewServeMux()
	apiUrl := env.GetOrElse("ASSET_API_URL", "https://api.binance.com")
	assetClient := binance.NewClient(apiUrl)

	setupRoutes(mux, assetClient)

	port := env.GetIntOrElse("PORT", DEFAULT_PORT)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	slog.Info("server listening...", "port", port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "error", err)
	}
}
