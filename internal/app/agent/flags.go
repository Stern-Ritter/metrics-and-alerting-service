package agent

import (
	"flag"

	"github.com/caarlos0/env"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// GetConfig initializes the agent config by parsing command-line flags and environment variables.
// It returns the initialized agent config and parsing error.
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
	flag.IntVar(&c.UpdateMetricsInterval, "p", 2, "interval for updating metrics in seconds")
	flag.IntVar(&c.SendMetricsInterval, "r", 10, "interval for sending metrics to the server in seconds")
	flag.IntVar(&c.RateLimit, "l", 1, "limit of concurrent requests to the server")
	flag.StringVar(&c.SecretKey, "k", "", "secret authentication key")
	flag.StringVar(&c.CryptoKeyPath, "crypto-key", "", "path to secret public key for asymmetric encryption")
	flag.Parse()
	err := utils.ValidateHostnamePort(c.SendMetricsURL)

	return err
}
