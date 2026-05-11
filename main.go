package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

var SHA = "dev"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(healthHandler))
	mux.Handle("/", http.FileServer(http.Dir("_site")))
	return mux
}

func main() {
	srv := &http.Server{
		Addr:              serverAddress(),
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", srv.Addr, "sha", SHA)
		serverErr <- srv.ListenAndServe()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server exited with error", "error", err)
			os.Exit(1)
		}
		slog.Info("server stopped")
		return
	case sig := <-signals:
		slog.Info("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	slog.Info("shutting down server", "timeout", shutdownTimeout.String())
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		if closeErr := srv.Close(); closeErr != nil {
			slog.Error("server close failed", "error", closeErr)
		}
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
