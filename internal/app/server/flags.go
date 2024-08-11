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
	URL             string `json:"address,omitempty"`
	StoreInterval   int    `json:"store_interval,omitempty"`
	FileStoragePath string `json:"store_file,omitempty"`
	Restore         bool   `json:"restore,omitempty"`
	DatabaseDSN     string `json:"database_dsn,omitempty"`
	GRPC            bool   `json:"grpc,omitempty"`
	SecretKey       string `json:"sign_key,omitempty"`
	CryptoKeyPath   string `json:"crypto_key,omitempty"`
	TLSCertPath     string `json:"tls_cert,omitempty"`
	TLSKeyPath      string `json:"tls_key,omitempty"`
	TrustedSubnet   string `json:"trusted_subnet,omitempty"`
	ShutdownTimeout int    `json:"shutdown_timeout,omitempty"`
	LoggerLvl       string `json:"logger_level,omitempty"`
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
	needParseJSONConfig := len(cfgFile) > 0
	if needParseJSONConfig {
		err = parseJSONConfig(&cfg, cfgFile)
		if err != nil {
			return cfg, err
		}
	}

	mergeDefaultConfig(&cfg, defaultCfg)

	trimStringVarsSpaces(&cfg)

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
	flag.BoolVar(&cfg.GRPC, "grpc", false, "grpc usage")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret authentication key")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to secret private key for asymmetric encryption")
	flag.StringVar(&cfg.TLSCertPath, "tls-cert", "", "path to tls certificate")
	flag.StringVar(&cfg.TLSKeyPath, "tls-key", "", "path to tls key")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "trusted subnet for agents")
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
	cfg.GRPC = utils.Coalesce(cfg.GRPC, jsonCfg.GRPC)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, jsonCfg.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, jsonCfg.CryptoKeyPath)
	cfg.TLSCertPath = utils.Coalesce(cfg.TLSCertPath, jsonCfg.TLSCertPath)
	cfg.TLSKeyPath = utils.Coalesce(cfg.TLSKeyPath, jsonCfg.TLSKeyPath)
	cfg.TrustedSubnet = utils.Coalesce(cfg.TrustedSubnet, jsonCfg.TrustedSubnet)
	cfg.ShutdownTimeout = utils.Coalesce(cfg.ShutdownTimeout, jsonCfg.ShutdownTimeout)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, jsonCfg.LoggerLvl)
}

func mergeDefaultConfig(cfg *config.ServerConfig, defaultCfg config.ServerConfig) {
	cfg.URL = utils.Coalesce(cfg.URL, defaultCfg.URL)
	cfg.StoreInterval = utils.Coalesce(cfg.StoreInterval, defaultCfg.StoreInterval)
	cfg.FileStoragePath = utils.Coalesce(cfg.FileStoragePath, defaultCfg.FileStoragePath)
	cfg.Restore = utils.Coalesce(cfg.Restore, defaultCfg.Restore)
	cfg.DatabaseDSN = utils.Coalesce(cfg.DatabaseDSN, defaultCfg.DatabaseDSN)
	cfg.GRPC = utils.Coalesce(cfg.GRPC, defaultCfg.GRPC)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, defaultCfg.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, defaultCfg.CryptoKeyPath)
	cfg.TLSCertPath = utils.Coalesce(cfg.TLSCertPath, defaultCfg.TLSCertPath)
	cfg.TLSKeyPath = utils.Coalesce(cfg.TLSKeyPath, defaultCfg.TLSKeyPath)
	cfg.TrustedSubnet = utils.Coalesce(cfg.TrustedSubnet, defaultCfg.TrustedSubnet)
	cfg.ConfigFile = utils.Coalesce(cfg.ConfigFile, defaultCfg.ConfigFile)
	cfg.ShutdownTimeout = utils.Coalesce(cfg.ShutdownTimeout, defaultCfg.ShutdownTimeout)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, defaultCfg.LoggerLvl)
}

func trimStringVarsSpaces(cfg *config.ServerConfig) {
	cfg.URL = strings.TrimSpace(cfg.URL)
	cfg.FileStoragePath = strings.TrimSpace(cfg.FileStoragePath)
	cfg.DatabaseDSN = strings.TrimSpace(cfg.DatabaseDSN)
	cfg.SecretKey = strings.TrimSpace(cfg.SecretKey)
	cfg.CryptoKeyPath = strings.TrimSpace(cfg.CryptoKeyPath)
	cfg.TLSCertPath = strings.TrimSpace(cfg.TLSCertPath)
	cfg.TLSKeyPath = strings.TrimSpace(cfg.TLSKeyPath)
	cfg.TrustedSubnet = strings.TrimSpace(cfg.TrustedSubnet)
	cfg.ConfigFile = strings.TrimSpace(cfg.ConfigFile)
	cfg.LoggerLvl = strings.TrimSpace(cfg.LoggerLvl)
}

func validateConfig(cfg config.ServerConfig) error {
	return utils.ValidateHostnamePort(cfg.URL)
}
