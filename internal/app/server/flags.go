package server

import (
	"flag"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/caarlos0/env"
)

func getConfig(c config.ServerConfig) (config.ServerConfig, error) {
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
	flag.Parse()
	err := utils.ValidateHostnamePort(c.URL)

	return err
}
