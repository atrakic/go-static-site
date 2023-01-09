package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format("02 Jan 2006 15:04:05 MST"))
}

func main() {
	fileHandler := http.StripPrefix("/", http.FileServer(http.Dir("_site/")))
	http.Handle("/", fileHandler)
	http.HandleFunc("/health", timeHandler)

	port := ":8080"
	fmt.Printf("Starting at %v\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
