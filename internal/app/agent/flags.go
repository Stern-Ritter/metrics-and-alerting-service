package agent

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type jsonConfig struct {
	SendMetricsURL        string `json:"address,omitempty"`
	SendMetricsEndPoint   string `json:"endpoint,omitempty"`
	UpdateMetricsInterval int    `json:"poll_interval,omitempty"`
	SendMetricsInterval   int    `json:"report_interval,omitempty"`
	MetricsBufferSize     int    `json:"metrics_buffer_size,omitempty"`
	RateLimit             int    `json:"rate_limit,omitempty"`
	GRPC                  bool   `json:"grpc,omitempty"`
	SecretKey             string `json:"sign_key,omitempty"`
	CryptoKeyPath         string `json:"crypto_key,omitempty"`
	TLSCertPath           string `json:"tls_cert,omitempty"`
	LoggerLvl             string `json:"logger_level,omitempty"`
}

// GetConfig initializes the agent config by parsing command-line flags, environment variables, and a JSON config file.
// It returns the initialized agent config and any parsing error encountered.
//
// The priority of configuration values is as follows (from highest to lowest):
// 1. Environment variables
// 2. Command-line flags
// 3. JSON config file
// 4. Default config
func GetConfig(defaultCfg config.AgentConfig) (config.AgentConfig, error) {
	cfg := config.AgentConfig{}

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

func parseFlags(cfg *config.AgentConfig) {
	flag.StringVar(&cfg.SendMetricsURL, "a", "", "address and port to run server in format <host>:<port>")
	flag.IntVar(&cfg.UpdateMetricsInterval, "p", 0, "interval for updating metrics in seconds")
	flag.IntVar(&cfg.SendMetricsInterval, "r", 0, "interval for sending metrics to the server in seconds")
	flag.IntVar(&cfg.RateLimit, "l", 0, "limit of concurrent requests to the server")
	flag.BoolVar(&cfg.GRPC, "grpc", false, "grpc usage")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret authentication key")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to secret public key for asymmetric encryption")
	flag.StringVar(&cfg.TLSCertPath, "tls-cert", "", "path to tls certificate")
	flag.StringVar(&cfg.ConfigFile, "c", "", "path to json config file")
	flag.Parse()
}

func parseJSONConfig(cfg *config.AgentConfig, fPath string) error {
	data, err := os.ReadFile(fPath)
	if err != nil {
		return fmt.Errorf("read config file %s: %w", fPath, err)
	}

	jsonCfg := jsonConfig{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	err = dec.Decode(&jsonCfg)
	if err != nil {
		return fmt.Errorf("parse config file %s: %w", fPath, err)
	}

	mergeJSONConfig(cfg, jsonCfg)

	return nil
}

func mergeJSONConfig(cfg *config.AgentConfig, jsonCfg jsonConfig) {
	cfg.SendMetricsURL = utils.Coalesce(cfg.SendMetricsURL, jsonCfg.SendMetricsURL)
	cfg.SendMetricsEndPoint = utils.Coalesce(cfg.SendMetricsEndPoint, jsonCfg.SendMetricsEndPoint)
	cfg.UpdateMetricsInterval = utils.Coalesce(cfg.UpdateMetricsInterval, jsonCfg.UpdateMetricsInterval)
	cfg.SendMetricsInterval = utils.Coalesce(cfg.SendMetricsInterval, jsonCfg.SendMetricsInterval)
	cfg.MetricsBufferSize = utils.Coalesce(cfg.MetricsBufferSize, jsonCfg.MetricsBufferSize)
	cfg.RateLimit = utils.Coalesce(cfg.RateLimit, jsonCfg.RateLimit)
	cfg.GRPC = utils.Coalesce(cfg.GRPC, jsonCfg.GRPC)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, jsonCfg.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, jsonCfg.CryptoKeyPath)
	cfg.TLSCertPath = utils.Coalesce(cfg.TLSCertPath, jsonCfg.TLSCertPath)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, jsonCfg.LoggerLvl)
}

func mergeDefaultConfig(cfg *config.AgentConfig, defaultCgf config.AgentConfig) {
	cfg.SendMetricsURL = utils.Coalesce(cfg.SendMetricsURL, defaultCgf.SendMetricsURL)
	cfg.SendMetricsEndPoint = utils.Coalesce(cfg.SendMetricsEndPoint, defaultCgf.SendMetricsEndPoint)
	cfg.UpdateMetricsInterval = utils.Coalesce(cfg.UpdateMetricsInterval, defaultCgf.UpdateMetricsInterval)
	cfg.SendMetricsInterval = utils.Coalesce(cfg.SendMetricsInterval, defaultCgf.SendMetricsInterval)
	cfg.MetricsBufferSize = utils.Coalesce(cfg.MetricsBufferSize, defaultCgf.MetricsBufferSize)
	cfg.RateLimit = utils.Coalesce(cfg.RateLimit, defaultCgf.RateLimit)
	cfg.GRPC = utils.Coalesce(cfg.GRPC, defaultCgf.GRPC)
	cfg.SecretKey = utils.Coalesce(cfg.SecretKey, defaultCgf.SecretKey)
	cfg.CryptoKeyPath = utils.Coalesce(cfg.CryptoKeyPath, defaultCgf.CryptoKeyPath)
	cfg.TLSCertPath = utils.Coalesce(cfg.TLSCertPath, defaultCgf.TLSCertPath)
	cfg.ConfigFile = utils.Coalesce(cfg.ConfigFile, defaultCgf.ConfigFile)
	cfg.LoggerLvl = utils.Coalesce(cfg.LoggerLvl, defaultCgf.LoggerLvl)
}

func trimStringVarsSpaces(cfg *config.AgentConfig) {
	cfg.SendMetricsURL = strings.TrimSpace(cfg.SendMetricsURL)
	cfg.SendMetricsEndPoint = strings.TrimSpace(cfg.SendMetricsEndPoint)
	cfg.SecretKey = strings.TrimSpace(cfg.SecretKey)
	cfg.CryptoKeyPath = strings.TrimSpace(cfg.CryptoKeyPath)
	cfg.TLSCertPath = strings.TrimSpace(cfg.TLSCertPath)
	cfg.ConfigFile = strings.TrimSpace(cfg.ConfigFile)
	cfg.LoggerLvl = strings.TrimSpace(cfg.LoggerLvl)
}

func validateConfig(cfg config.AgentConfig) error {
	return utils.ValidateHostnamePort(cfg.SendMetricsURL)
}
