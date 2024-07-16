package agent

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gentleman.v2/context"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
)

type MockHandler struct {
	mock.Mock
}

func (h *MockHandler) Next(ctx *context.Context) {
	h.Called(ctx)
}

func (h *MockHandler) Stop(ctx *context.Context) {
	h.Called(ctx)
}

func (h *MockHandler) Error(ctx *context.Context, err error) {
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
			ctx := &context.Context{
				Request: req,
			}

			agent := &Agent{
				Config: &config.AgentConfig{
					SecretKey: tt.agentSignKey,
				},
			}

			mockHandler := &MockHandler{}
			mockHandler.On("Next", ctx).Once()

			agent.SignMiddleware(ctx, mockHandler)

			gotSign := req.Header.Get(signKey)
			assert.Equal(t, tt.expectedSign, gotSign, "request body sign should be: %s, got: %s", tt.expectedSign, gotSign)
			mockHandler.AssertExpectations(t)
		})
	}
}
