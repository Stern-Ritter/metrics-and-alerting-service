package agent

import (
	"context"
	"crypto/rsa"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/h2non/gentleman.v2"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/agent"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	crypto "github.com/Stern-Ritter/metrics-and-alerting-service/internal/crypto/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// Count of agent setting up and managing tasks
const (
	taskCount = 3
)

// Run starts the agent, setting up and managing tasks.
// It returns an error if there are issues starting the agent.
func Run(config *config.AgentConfig, logger *logger.AgentLogger) error {
	httpClient := gentleman.New()
	cache := storage.NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics, logger)
	runtimeMonitor := monitors.RuntimeMonitor{}
	utilMonitor := monitors.UtilMonitor{}
	random := utils.NewRandom()

	var rsaPublicKey *rsa.PublicKey

	isEncryptionEnabled := len(strings.TrimSpace(config.CryptoKeyPath)) != 0
	if isEncryptionEnabled {
		key, err := crypto.GetRSAPublicKey(config.CryptoKeyPath)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "get rsa public key"))
			return err
		}
		rsaPublicKey = key
	}

	agent := service.NewAgent(httpClient, &cache, &runtimeMonitor, &utilMonitor, &random, config, rsaPublicKey, logger)

	agent.HTTPClient.URL(agent.Config.SendMetricsURL)
	agent.HTTPClient.UseHandler("before dial", compress.GzipMiddleware)
	agent.HTTPClient.UseHandler("before dial", agent.EncryptMiddleware)
	agent.HTTPClient.UseHandler("before dial", agent.SignMiddleware)

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent.StartSendMetricsWorkerPool()

	service.SetInterval(ctx, &wg, agent.UpdateRuntimeMetrics, time.Duration(agent.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, agent.UpdateUtilMetrics, time.Duration(agent.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, agent.SendMetrics, time.Duration(agent.Config.SendMetricsInterval)*time.Second)

	wg.Wait()
	agent.StopTasks()

	return nil
}
