package server

import (
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

type Server struct {
	Storage storage.ServerStorage
	Config  *config.ServerConfig
	Logger  *logger.ServerLogger
}

func NewServer(storage storage.ServerStorage, config *config.ServerConfig,
	logger *logger.ServerLogger) *Server {
	return &Server{Storage: storage, Config: config, Logger: logger}
}
