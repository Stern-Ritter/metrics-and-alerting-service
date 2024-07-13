package server

import (
	"crypto/rsa"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
)

// Server is the server for handling requests to work with metrics representation.
type Server struct {
	MetricService *MetricService       // MetricService handles requests to work with metrics representation
	Config        *config.ServerConfig // Config holds the server configuration.
	rsaPrivateKey *rsa.PrivateKey      // rsaPrivateKey is secret private key for asymmetric encryption
	Logger        *logger.ServerLogger // Logger is used for logging server events.
}

// NewServer is constructor for creating a new Server instance.
func NewServer(metricService *MetricService, config *config.ServerConfig, rsaPrivateKey *rsa.PrivateKey,
	logger *logger.ServerLogger) *Server {
	return &Server{
		MetricService: metricService,
		Config:        config,
		rsaPrivateKey: rsaPrivateKey,
		Logger:        logger,
	}
}
