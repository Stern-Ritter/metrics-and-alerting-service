package server

import (
	"crypto/rsa"
	"net"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics"
)

// Server is the server for handling requests to work with metrics representation.
type Server struct {
	MetricService *MetricService       // MetricService handles requests to work with metrics representation
	Config        *config.ServerConfig // Config holds the server configuration
	rsaPrivateKey *rsa.PrivateKey      // rsaPrivateKey is secret private key for asymmetric encryption
	trustedSubnet *net.IPNet           // trustedSubnet is trusted subnet for agents
	Logger        *logger.ServerLogger // Logger is used for logging server events
	pb.UnimplementedMetricsServer
}

// NewServer is constructor for creating a new Server instance.
func NewServer(metricService *MetricService, config *config.ServerConfig, rsaPrivateKey *rsa.PrivateKey,
	trustedSubnet *net.IPNet, logger *logger.ServerLogger) *Server {
	return &Server{
		MetricService: metricService,
		Config:        config,
		rsaPrivateKey: rsaPrivateKey,
		trustedSubnet: trustedSubnet,
		Logger:        logger,
	}
}
