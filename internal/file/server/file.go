package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"go.uber.org/zap"
)

type FileStorage interface {
	Load() error
	Save() error
	Close() error
	SetSaveInterval(interval int)
	FileStorageMiddleware(next http.Handler) http.Handler
}

type Metrics struct {
	Gauges   map[string]metrics.GaugeMetric   `json:"gauges"`
	Counters map[string]metrics.CounterMetric `json:"counters"`
}

type ServerFileStorage struct {
	file *os.File

	encoder *json.Encoder
	decoder *json.Decoder

	synchronous bool
	data        Metrics

	storage storage.ServerStorage
}

func NewServerFileStorage(fname string, storage storage.ServerStorage) (*ServerFileStorage, error) {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &ServerFileStorage{
		file:        file,
		encoder:     json.NewEncoder(file),
		decoder:     json.NewDecoder(file),
		storage:     storage,
		synchronous: true,
		data: Metrics{
			Gauges:   make(map[string]metrics.GaugeMetric),
			Counters: make(map[string]metrics.CounterMetric),
		},
	}, nil
}

func (s *ServerFileStorage) Load() error {
	err := s.decoder.Decode(&s.data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	if len(s.data.Gauges) != 0 {
		s.storage.SetGaugeMetircs(s.data.Gauges)
	}
	if len(s.data.Counters) != 0 {
		s.storage.SetCounterMetrics(s.data.Counters)
	}

	return nil
}

func (s *ServerFileStorage) Save() error {
	gauges, counters := s.storage.GetMetrics()
	s.data.Gauges = gauges
	s.data.Counters = counters

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = s.encoder.Encode(s.data)
	return err
}

func (s *ServerFileStorage) Close() error {
	err := s.Save()
	if err != nil {
		logger.Log.Error(err.Error(), zap.String("event", "save to file storage when closing file storage"))
		return err
	}
	err = s.file.Close()
	logger.Log.Info("Close file storage", zap.String("event", "close file storage"))
	return err
}

func (s *ServerFileStorage) SetSaveInterval(interval int) {
	if interval <= 0 {
		return
	}

	s.synchronous = false
	logger.Log.Info("Start async save to file storage", zap.String("event", "start async save to file storage"), zap.Int("interval", interval))

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for range ticker.C {
			if err := s.Save(); err != nil {
				logger.Log.Error(err.Error(), zap.String("event", "async save to file storage"))
			} else {
				logger.Log.Info("Success async save to file storage", zap.String("event", "async save to file storage"))
			}
		}
	}()

}

func (s *ServerFileStorage) FileStorageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if s.synchronous {
			err := s.Save()
			if err != nil {
				logger.Log.Error(err.Error(), zap.String("event", "sync save to file storage"))
			} else {
				logger.Log.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
			}
		}
	})
}