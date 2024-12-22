package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var SHA string

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format("02 Jan 2006 15:04:05 MST"))
	fmt.Fprint(w, "\nSHA: ", SHA)
}

func main() {
	fileHandler := http.StripPrefix("/", http.FileServer(http.Dir("_site/")))
	http.Handle("/", fileHandler)
	http.HandleFunc("/health", timeHandler)

	port := ":80"
	slog.Info("Starting server", "addr", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error(err.Error())
	}
}
