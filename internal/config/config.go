package config

import "time"

type AgentConfig struct {
	SendMetricsURL        URL
	SendMetricsEndPoint   string
	UpdateMetricsInterval time.Duration
	SendMetricsInterval   time.Duration
}

type ServerConfig struct {
	URL URL
}
