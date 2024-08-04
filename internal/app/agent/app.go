package agent

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"gopkg.in/h2non/gentleman.v2"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/agent"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics"
)

// Run starts the agent, setting up and managing tasks.
// It returns an error if there are issues starting the agent.
func Run(config *config.AgentConfig, logger *logger.AgentLogger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cache := storage.NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics, logger)
	runtimeMonitor := monitors.RuntimeMonitor{}
	utilMonitor := monitors.UtilMonitor{}
	random := utils.NewRandom()

	rsaPublicKey, err := service.GetRSAPublicKey(config.CryptoKeyPath)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "get rsa public key"))
	}

	agent := service.NewAgent(&cache, &runtimeMonitor, &utilMonitor, &random, config, rsaPublicKey, logger)

	tasksWg := sync.WaitGroup{}
	service.SetInterval(ctx, &tasksWg, agent.UpdateRuntimeMetrics, time.Duration(agent.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &tasksWg, agent.UpdateUtilMetrics, time.Duration(agent.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &tasksWg, agent.SendMetrics, time.Duration(agent.Config.SendMetricsInterval)*time.Second)

	workersWg := sync.WaitGroup{}

	if config.GRPC {
		opts := make([]grpc.DialOption, 0)

		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
		isEncryptionEnabled := len(agent.Config.TLSCertPath) > 0
		if isEncryptionEnabled {
			creds, err := credentials.NewClientTLSFromFile(agent.Config.TLSCertPath, "")
			if err != nil {
				logger.Fatal(err.Error(), zap.String("event", "load credentials"))
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		opts = append(opts, grpc.WithUnaryInterceptor(agent.SignInterceptor))

		conn, err := grpc.Dial(config.SendMetricsURL, opts...)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "get grpc connection"))
		}
		defer conn.Close()
		client := pb.NewMetricsClient(conn)
		agent.SetGRPCClient(client)

		agent.StartSendMetricsWorkerPool(&workersWg, agent.SendMetricsWithGrpcWorker)
	} else {
		client := gentleman.New()
		client.URL(config.SendMetricsURL)
		client.UseHandler("before dial", compress.GzipMiddleware)
		client.UseHandler("before dial", agent.EncryptMiddleware)
		client.UseHandler("before dial", agent.SignMiddleware)
		agent.SetHTTPClient(client)

		agent.StartSendMetricsWorkerPool(&workersWg, agent.SendMetricsWithHTTPWorker)
	}

	tasksWg.Wait()
	agent.StopSendMetricsWorkerPool()
	workersWg.Wait()

	return nil
}
