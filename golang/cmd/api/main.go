package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fernandezpablo85/sdi/internal/env"
)

const DEFAULT_PORT = 8080

func get(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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

func assetHanlder(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		slog.Debug("missing asset name", "status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing asset name"))
		return
	}
	slog.Info("fetching asset price", "asset", name)
	fmt.Fprintf(w, "Fetching %s price...\n", name)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", get(healthzHandler))
	mux.HandleFunc("/v1/asset", get(assetHanlder))

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
