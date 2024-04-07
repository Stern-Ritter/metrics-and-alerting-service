package agent

import (
	"go.uber.org/zap"
)

func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}
