package main

import (
	"log"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

func main() {
	storage := storage.NewServerMemStorage()
	server := service.NewServer(&storage)

	err := app.Run(server)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
