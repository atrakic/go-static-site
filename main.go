package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

var SHA = "dev"

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

type healthResponse struct {
	Time string `json:"time"`
	SHA  string `json:"sha"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	default:
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := json.NewEncoder(w).Encode(healthResponse{
		Time: time.Now().UTC().Format(time.RFC3339),
		SHA:  SHA,
	}); err != nil {
		slog.Error("write health response", "error", err)
	}
}

func serverAddress() string {
	port := envOrDefault("PORT", "8080")
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func staticDir() string {
	return envOrDefault("HTTP_DIR", "_site")
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(healthHandler))
	mux.Handle("/", http.FileServer(http.Dir(staticDir())))
	return mux
}

func main() {
	addr := serverAddress()
	slog.Info("starting server", "addr", addr, "sha", SHA)
	if err := http.ListenAndServe(addr, newMux()); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
