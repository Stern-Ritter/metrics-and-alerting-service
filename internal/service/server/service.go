package server

import (
	file "github.com/Stern-Ritter/metrics-and-alerting-service/internal/file/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

type Server struct {
	storage     storage.ServerStorage
	FileStorage *file.FileStorage
}

func NewServer(storage storage.ServerStorage) *Server {
	return &Server{storage: storage}
}

func (s *Server) AddFileStorage(fname string) error {
	fileStorage, err := file.NewFileStorage(fname, s.storage)
	if err != nil {
		return err
	}

	s.FileStorage = fileStorage
	return nil
}
