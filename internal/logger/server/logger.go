package server

import (
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	interceptor "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/grpc"
)

// ServerLogger wraps a zap.Logger.
type ServerLogger struct {
	*zap.Logger
}

// Initialize initializes a ServerLogger with the specified logging level.
func Initialize(level string) (*ServerLogger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ServerLogger{logger}, nil
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write writes the data to response body.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader write response headers and set response status code.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LoggerMiddleware is an HTTP middleware fot logging the request details:
// request uri, request method, process duration, response status, response body size
func (logger *ServerLogger) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logger.Info(
			"Request received: ",
			zap.String("event", "request"),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}

// LoggerInterceptor returns a new UnaryServerInterceptor that logs the details of each gRPC call.
// It uses the provided logger to log events related to the gRPC calls. The interceptor logs
// only when the call finishes.
func LoggerInterceptor(logger interceptor.Logger) grpc.UnaryServerInterceptor {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
	}
	return logging.UnaryServerInterceptor(interceptor.NewInterceptorLogger(logger), opts...)
}
