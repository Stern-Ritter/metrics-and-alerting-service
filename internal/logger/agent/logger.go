package agent

import (
	"go.uber.org/zap"
)

// AgentLogger wraps a zap.Logger.
type AgentLogger struct {
	*zap.Logger
}

// Initialize initializes a zap.Logger with the specified logging level.
func Initialize(level string) (*AgentLogger, error) {
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

	return &AgentLogger{log}, nil
}
