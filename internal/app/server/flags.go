package server

import (
	"flag"

	"github.com/caarlos0/env"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// GetConfig initializes the server config by parsing command-line flags and environment variables.
// It returns the initialized server config and parsing error.
func GetConfig(c config.ServerConfig) (config.ServerConfig, error) {
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

func parseFlags(c *config.ServerConfig) error {
	flag.StringVar(&c.URL, "a", "localhost:8080", "address and port to run server in format <host>:<port>")
	flag.IntVar(&c.StoreInterval, "i", 300, "interval to store metrics to file in seconds")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "metrics storage file path")
	flag.BoolVar(&c.Restore, "r", true, "will metrics be restored from the file")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&c.SecretKey, "k", "", "secret authentication key")
	flag.Parse()
	err := utils.ValidateHostnamePort(c.URL)

	return err
}
