package main

import (
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/handlers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetricHandler)
	err := http.ListenAndServe(`:8080`, mux)
	return err
}
