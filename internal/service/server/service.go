package server

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

type Server struct {
	storage storage.ServerStorage
}

func NewServer(storage storage.ServerStorage) *Server {
	return &Server{storage}
}
