package agent

import (
	"go.uber.org/zap"
)

// Initialize initializes a zap.Logger with the specified logging level.
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
