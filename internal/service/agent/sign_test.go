package agent

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gcontext "gopkg.in/h2non/gentleman.v2/context"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

type SignMockHandler struct {
	mock.Mock
}

func (h *SignMockHandler) Next(ctx *gcontext.Context) {
	h.Called(ctx)
}

func (h *SignMockHandler) Stop(ctx *gcontext.Context) {
	h.Called(ctx)
}

func (h *SignMockHandler) Error(ctx *gcontext.Context, err error) {
	h.Called(ctx, err)
}

func TestSignMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		agentSignKey string
		requestBody  string
		expectedSign string
	}{
		{
			name:         "should sign request body when sign key is defined and request body isn`t empty",
			agentSignKey: "secret",
			requestBody:  "The Ultimate Question of Life, the Universe, and Everything",
			expectedSign: "d447119d670edac6cf426ba1d905508636f43a852ddae829a7062eae58ab845a",
		},
		{
			name:         "shouldn`t sign request body when sign key is defined and request body is empty",
			agentSignKey: "secret",
			requestBody:  "",
			expectedSign: "",
		},
		{
			name:         "shouldn`t sign request body when sign key isn`t defined and request body isn`t empty",
			agentSignKey: "",
			requestBody:  "The Ultimate Question of Life, the Universe, and Everything",
			expectedSign: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewReader([]byte(tt.requestBody))),
			}
			ctx := &gcontext.Context{
				Request: req,
			}

			agent := &Agent{
				Config: &config.AgentConfig{
					SecretKey: tt.agentSignKey,
				},
			}

			mockHandler := &SignMockHandler{}
			mockHandler.On("Next", ctx).Once()

			agent.SignMiddleware(ctx, mockHandler)

			gotSign := req.Header.Get(signKey)
			assert.Equal(t, tt.expectedSign, gotSign, "request body sign should be: %s, got: %s", tt.expectedSign, gotSign)
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestSignInterceptor(t *testing.T) {
	secretKey := "secret-key"

	tests := []struct {
		name        string
		secretKey   string
		req         interface{}
		expectedErr error
		expectedMD  metadata.MD
	}{
		{
			name:      "Valid request",
			secretKey: secretKey,
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: nil,
			expectedMD: metadata.MD{
				"hashsha256": {"a3fe670ac7a7e88578e5f20b5a1619c217da021335580c10c419cbf8472d9bd9"},
			},
		},
		{
			name:        "Request not proto.Message",
			secretKey:   secretKey,
			req:         "invalid request",
			expectedErr: status.Errorf(codes.Internal, "sign interceptor: request isn't a proto.Message"),
			expectedMD:  metadata.MD{},
		},
		{
			name:      "No secret key in agent config",
			secretKey: "",
			req: &pb.MetricsV1ServiceUpdateMetricRequest{Metric: &pb.MetricData{
				Name:        "Alloc",
				Type:        "gauge",
				MetricValue: &pb.MetricData_Value{Value: 22.22},
			}},
			expectedErr: nil,
			expectedMD:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{Config: &config.AgentConfig{SecretKey: tt.secretKey}}

			ctx := context.Background()

			testInvoker := func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				md, _ := metadata.FromOutgoingContext(ctx)
				if tt.expectedErr == nil {
					assert.Equal(t, tt.expectedMD, md, "metadata does not match")
				}

				return nil
			}

			err := agent.SignInterceptor(ctx, "method", tt.req, nil, nil, testInvoker)

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err, "should return error: %s, but got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn`t return error, but got: %s", err)
			}
		})
	}
}
