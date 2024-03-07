package config

import "time"

type AgentConfig struct {
	UpdateMetricsInterval time.Duration
	SendMetricsInterval   time.Duration
	SendMetricsEndPoint   string
}

type ServerConfig struct {
	URL string
}

var MonitoringAgentConfig = AgentConfig{
	UpdateMetricsInterval: 2 * time.Second,
	SendMetricsInterval:   10 * time.Second,
	SendMetricsEndPoint:   "http://localhost:8080/update",
}

var MetricsServerConfig = ServerConfig{
	URL: `:8080`,
}
