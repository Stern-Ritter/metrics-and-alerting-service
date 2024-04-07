package server

import (
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
)

type Server struct {
	MetricService *MetricService
	Config        *config.ServerConfig
	Logger        *logger.ServerLogger
}

func NewServer(metricService *MetricService, config *config.ServerConfig,
	logger *logger.ServerLogger) *Server {
	return &Server{MetricService: metricService, Config: config, Logger: logger}
}
