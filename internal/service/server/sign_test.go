package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

func TestSignMiddleware(t *testing.T) {
	type request struct {
		sign string
		body string
	}

	type response struct {
		status int
		sign   string
		body   string
	}

	tests := []struct {
		name          string
		serverSignKey string
		request       request
		response      response
	}{
		{
			name:          "should return status ok and signed body when sign key is defined and request body has valid sign",
			serverSignKey: "secret",
			request: request{
				sign: "d447119d670edac6cf426ba1d905508636f43a852ddae829a7062eae58ab845a",
				body: "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusOK,
				sign:   "1EcRnWcO2sbPQmuh2QVQhjb0OoUt2ugppwYurlirhFo=",
				body:   "The Ultimate Question of Life, the Universe, and Everything",
			},
		},
		{
			name:          "should return status bad request when sign key is defined and request body has invalid sign",
			serverSignKey: "secret",
			request: request{
				sign: "invalid sign",
				body: "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusBadRequest,
				sign:   "",
				body:   "Invalid request body sign\n",
			},
		},
		{
			name:          "should return status ok and unsigned body when sign key isn`t defined and request body has valid sign",
			serverSignKey: "",
			request: request{
				sign: "d447119d670edac6cf426ba1d905508636f43a852ddae829a7062eae58ab845a",
				body: "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusOK,
				sign:   "",
				body:   "The Ultimate Question of Life, the Universe, and Everything",
			},
		},
		{
			name:          "should return status ok and unsigned body when sign key isn`t defined and request body has invalid sign",
			serverSignKey: "",
			request: request{
				sign: "invalid sign",
				body: "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusOK,
				sign:   "",
				body:   "The Ultimate Question of Life, the Universe, and Everything",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte(tt.request.body))))
			req.Header.Set(signKey, tt.request.sign)

			server := &Server{
				Config: &config.ServerConfig{
					SecretKey: tt.serverSignKey,
				},
			}

			handler := server.SignMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				_, _ = w.Write(body)
			}))

			r := httptest.NewRecorder()
			handler.ServeHTTP(r, req)

			resp := r.Result()
			data, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "unexpected error when read response body")
			defer resp.Body.Close()

			gotStatus := resp.StatusCode
			gotSign := resp.Header.Get(signKey)
			gotBody := string(data)
			assert.Equal(t, tt.response.status, gotStatus, "response status code should be %d, got %d",
				tt.response.status, gotStatus)
			assert.Equal(t, tt.response.sign, gotSign, "response body sign should be: %s, got: %s",
				tt.response.sign, gotSign)
			assert.Equal(t, tt.response.body, gotBody, "response body should be: %s, got: %s",
				tt.response.body, gotBody)
		})
	}
}

func TestSignInterceptor(t *testing.T) {
	secretKey := "secret-key"

	tests := []struct {
		name        string
		secretKey   string
		md          metadata.MD
		req         interface{}
		expectedErr error
	}{
		{
			name:      "Valid sign in request metadata",
			secretKey: secretKey,
			md:        metadata.Pairs(signKey, "a3fe670ac7a7e88578e5f20b5a1619c217da021335580c10c419cbf8472d9bd9"),
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: nil,
		},
		{
			name:      "No request metadata",
			secretKey: secretKey,
			md:        nil,
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: status.Errorf(codes.InvalidArgument, "sign interceptor: missing request metadata"),
		},
		{
			name:      "No sign in request metadata",
			secretKey: secretKey,
			md:        metadata.Pairs(),
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: status.Errorf(codes.InvalidArgument, "sign interceptor: missing request sign"),
		},
		{
			name:        "Request not proto.Message",
			secretKey:   secretKey,
			md:          metadata.Pairs(signKey, "a3fe670ac7a7e88578e5f20b5a1619c217da021335580c10c419cbf8472d9bd9"),
			req:         "invalid request",
			expectedErr: status.Errorf(codes.Internal, "sign interceptor: request isn't a proto.Message"),
		},
		{
			name:      "Invalid sign in request metadata",
			secretKey: secretKey,
			md:        metadata.Pairs(signKey, "invalid"),
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: status.Errorf(codes.Unauthenticated, "sign interceptor: invalid sign"),
		},
		{
			name:      "No secret key in server config",
			secretKey: "",
			md:        nil,
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{Config: &config.ServerConfig{SecretKey: tt.secretKey}}

			ctx := context.Background()
			if tt.md != nil {
				ctx = metadata.NewIncomingContext(ctx, tt.md)
			}

			_, err := server.SignInterceptor(ctx, tt.req, nil,
				func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "should not return error, but got: %s", err)
			}
		})
	}
}
