package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type jsonConfig struct {
	URL             string `json:"address"`
	StoreInterval   int    `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	Restore         bool   `json:"restore"`
	DatabaseDSN     string `json:"database_dsn"`
	SecretKey       string `json:"sign_key"`
	CryptoKeyPath   string `json:"crypto_key"`
	LoggerLvl       string `json:"logger_level"`
}

// GetConfig initializes the server config by parsing command-line flags, environment variables, and a JSON config file.
// It returns the initialized server config and any parsing error encountered.
//
// The priority of configuration values is as follows (from highest to lowest):
// 1. Environment variables
// 2. Command-line flags
// 3. JSON config file
// 4. Default config
func GetConfig(defaultCfg config.ServerConfig) (config.ServerConfig, error) {
	cfg := config.ServerConfig{}

	parseFlags(&cfg)

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	cfgFile := strings.TrimSpace(cfg.ConfigFile)
	if len(cfgFile) > 0 {
		err = parseJSONConfig(&cfg, cfgFile)
		if err != nil {
			return cfg, err
		}
	}

	mergeDefaultConfig(&cfg, defaultCfg)

	err = validateConfig(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func parseFlags(cfg *config.ServerConfig) {
	flag.StringVar(&cfg.URL, "a", "", "address and port to run server in format <host>:<port>")
	flag.IntVar(&cfg.StoreInterval, "i", 0, "interval to store metrics to file in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "metrics storage file path")
	flag.BoolVar(&cfg.Restore, "r", false, "will metrics be restored from the file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret authentication key")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to secret private key for asymmetric encryption")
	flag.StringVar(&cfg.ConfigFile, "c", "", "path to json config file")
	flag.Parse()
}

func parseJSONConfig(cfg *config.ServerConfig, fPath string) error {
	data, err := os.ReadFile(fPath)
	if err != nil {
		return fmt.Errorf("read config file %s: %w", fPath, err)
	}

	jsonCfg := jsonConfig{}
	err = json.Unmarshal(data, &jsonCfg)
	if err != nil {
		return fmt.Errorf("parse config file %s: %w", fPath, err)
	}

	mergeJSONConfig(cfg, jsonCfg)
	return nil
}

func mergeJSONConfig(cfg *config.ServerConfig, jsonCfg jsonConfig) {
	cfg.URL = utils.Coalesce(cfg.URL, jsonCfg.URL)
	cfg.StoreInterval = utils.Coalesce(cfg.StoreInterval, jsonCfg.StoreInterval)
	cfg.FileStoragePath = utils.Coalesce(cfg.FileStoragePath, jsonCfg.FileStoragePath)
	cfg.Restore = utils.Coalesce(cfg.Restore, jsonCfg.Restore)
	cfg.DatabaseDSN = utils.Coalesce(cfg.DatabaseDSN, jsonCfg.DatabaseDSN)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, jsonCfg.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, jsonCfg.CryptoKeyPath)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, jsonCfg.LoggerLvl)
}

func mergeDefaultConfig(cfg *config.ServerConfig, defaultCfg config.ServerConfig) {
	cfg.URL = utils.Coalesce(cfg.URL, defaultCfg.URL)
	cfg.StoreInterval = utils.Coalesce(cfg.StoreInterval, defaultCfg.StoreInterval)
	cfg.FileStoragePath = utils.Coalesce(cfg.FileStoragePath, defaultCfg.FileStoragePath)
	cfg.Restore = utils.Coalesce(cfg.Restore, defaultCfg.Restore)
	cfg.DatabaseDSN = utils.Coalesce(cfg.DatabaseDSN, defaultCfg.DatabaseDSN)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, defaultCfg.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, defaultCfg.CryptoKeyPath)
	cfg.ConfigFile = utils.Coalesce(cfg.ConfigFile, defaultCfg.ConfigFile)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, defaultCfg.LoggerLvl)
}

func validateConfig(cfg config.ServerConfig) error {
	return utils.ValidateHostnamePort(cfg.URL)
}
