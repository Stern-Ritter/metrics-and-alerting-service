package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"go.uber.org/zap"
)

type StorageState struct {
	Gauges   map[string]metrics.GaugeMetric   `json:"gauges"`
	Counters map[string]metrics.CounterMetric `json:"counters"`
}

type ServerStorage interface {
	Storage
	Restore(fname string) error
	Save(fname string) error
	SetSaveInterval(fname string, interval int)
}

type ServerMemStorage struct {
	MemStorage
	Logger *logger.ServerLogger
}

func NewServerMemStorage(logger *logger.ServerLogger) ServerMemStorage {
	return ServerMemStorage{
		MemStorage: MemStorage{
			gauges:   make(map[string]metrics.GaugeMetric),
			counters: make(map[string]metrics.CounterMetric),
		},
		Logger: logger,
	}
}

func (s *ServerMemStorage) Restore(fname string) error {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	state := StorageState{
		Gauges:   make(map[string]metrics.GaugeMetric),
		Counters: make(map[string]metrics.CounterMetric),
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil
	}
	data := scanner.Bytes()

	err = json.Unmarshal(data, &state)
	if err != nil {
		return err
	}

	s.gaugesMu.Lock()
	s.gauges = state.Gauges
	s.gaugesMu.Unlock()

	s.countersMu.Lock()
	s.counters = state.Counters
	s.countersMu.Unlock()

	return nil
}

func (s *ServerMemStorage) Save(fname string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	s.gaugesMu.Lock()
	s.countersMu.Lock()
	state := StorageState{
		Gauges:   s.gauges,
		Counters: s.counters,
	}

	data, err := json.Marshal(&state)
	if err != nil {
		return err
	}
	s.gaugesMu.Unlock()
	s.countersMu.Unlock()

	_, err = file.Write(data)
	return err
}

func (s *ServerMemStorage) SetSaveInterval(fname string, interval int) {
	if interval <= 0 {
		return
	}

	s.Logger.Info("Start async save to file storage", zap.String("event", "start async save to file storage"),
		zap.String("file name", fname), zap.Int("interval", interval))
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for range ticker.C {
			if err := s.Save(fname); err != nil {
				s.Logger.Error(err.Error(), zap.String("event", "async save to file storage"))
			} else {
				s.Logger.Info("Success async save to file storage", zap.String("event", "async save to file storage"))
			}
		}
	}()
}
