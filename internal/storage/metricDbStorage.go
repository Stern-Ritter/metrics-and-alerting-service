package storage

import (
	"database/sql"

	"context"
)

type DBStorage interface {
	Ping(ctx context.Context) error
}

type MetricStore struct {
	db *sql.DB
}

func NewMetricStore(db *sql.DB) MetricStore {
	return MetricStore{db: db}
}

func (s *MetricStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
