package main

import (
	"log"
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	handlers "github.com/Stern-Ritter/metrics-and-alerting-service/internal/transport"
	"github.com/go-chi/chi"
)

func main() {
	config := config.ServerConfig{}
	err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewServerMemStorage()
	run(&storage, config)
}

func run(storage *storage.ServerMemStorage, config config.ServerConfig) {
	router := chi.NewRouter()
	router.Get("/", handlers.GetMetricsHandler(storage))
	router.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(storage))
	router.Get("/value/{type}/{name}", handlers.GetMetricHandler(storage))

	err := http.ListenAndServe(config.URL, router)
	log.Fatal(err)
}
