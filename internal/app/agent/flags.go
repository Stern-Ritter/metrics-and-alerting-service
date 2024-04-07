package agent

import (
	"flag"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/caarlos0/env"
)

func GetConfig(c config.AgentConfig) (config.AgentConfig, error) {
	err := parseFlags(&c)
	if err != nil {
		return c, err
	}

	err = env.Parse(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func parseFlags(c *config.AgentConfig) error {
	flag.StringVar(&c.SendMetricsURL, "a", "localhost:8080", "address and port to run server in format <host>:<port>")
	flag.IntVar(&c.UpdateMetricsInterval, "p", 2, "interval to update metrics in seconds")
	flag.IntVar(&c.SendMetricsInterval, "r", 10, "interval for sending metrics to the server in seconds")
	flag.Parse()
	err := utils.ValidateHostnamePort(c.SendMetricsURL)

	return err
}
