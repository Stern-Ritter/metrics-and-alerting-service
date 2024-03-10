package config

import (
	"flag"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type AgentConfig struct {
	SendMetricsURL        string `env:"ADDRESS"`
	SendMetricsEndPoint   string
	UpdateMetricsInterval int `env:"POLL_INTERVAL"`
	SendMetricsInterval   int `env:"REPORT_INTERVAL"`
}

func (c *AgentConfig) ParseFlags() error {
	flag.StringVar(&c.SendMetricsURL, "a", "localhost:8080", "address and port to run server in format <host>:<port>")
	flag.IntVar(&c.UpdateMetricsInterval, "p", 2, "interval to update metrics in seconds")
	flag.IntVar(&c.SendMetricsInterval, "r", 10, "interval for sending metrics to the server in seconds")
	flag.Parse()
	err := utils.ValidateHostnamePort(c.SendMetricsURL)

	return err
}

type ServerConfig struct {
	URL string `env:"ADDRESS"`
}

func (c *ServerConfig) ParseFlags() error {
	flag.StringVar(&c.URL, "a", "localhost:8080", "address and port to run server in format <host>:<port>")
	flag.Parse()
	err := utils.ValidateHostnamePort(c.URL)

	return err
}
