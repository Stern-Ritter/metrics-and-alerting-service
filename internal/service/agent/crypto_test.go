package agent

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gentleman.v2/context"
)

type EncryptMockHandler struct {
	mock.Mock
}

func (h *EncryptMockHandler) Next(ctx *context.Context) {
	h.Called(ctx)
}

func (h *EncryptMockHandler) Stop(ctx *context.Context) {
	h.Called(ctx)
}

func (h *EncryptMockHandler) Error(ctx *context.Context, err error) {
	h.Called(ctx, err)
}

func TestEncryptMiddleware(t *testing.T) {
	tests := []struct {
		name                   string
		hasCryptoKey           bool
		requestBody            string
		isRequestBodyEncrypted bool
	}{
		{
			name:                   "should encrypt request body when crypto key is defined and request body isn`t empty",
			hasCryptoKey:           true,
			requestBody:            "The Ultimate Question of Life, the Universe, and Everything",
			isRequestBodyEncrypted: true,
		},
		{
			name:                   "shouldn`t encrypt request body when crypto key is defined and request body is empty",
			hasCryptoKey:           true,
			requestBody:            "",
			isRequestBodyEncrypted: false,
		},
		{
			name:                   "shouldn`t encrypt request body when crypto key isn`t defined and request body isn`t empty",
			hasCryptoKey:           false,
			requestBody:            "The Ultimate Question of Life, the Universe, and Everything",
			isRequestBodyEncrypted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err, "unexpected error when generate private rsa key")
			publicKey := &privateKey.PublicKey

			req := &http.Request{
				Header: http.Header{},
				Body:   io.NopCloser(bytes.NewReader([]byte(tt.requestBody))),
			}
			ctx := &context.Context{
				Request: req,
			}

			var rsaPublicKey *rsa.PublicKey
			if tt.hasCryptoKey {
				rsaPublicKey = publicKey
			}

			agent := &Agent{
				rsaPublicKey: rsaPublicKey,
			}

			mockHandler := &EncryptMockHandler{}
			mockHandler.On("Next", ctx).Once()

			agent.EncryptMiddleware(ctx, mockHandler)

			body, err := io.ReadAll(ctx.Request.Body)
			require.NoError(t, err, "unexpected error when read request body")

			var expectedBody string
			if tt.isRequestBodyEncrypted {
				decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, body)
				require.NoError(t, err, "unexpected error when decrypt request body")
				expectedBody = string(decryptedBody)
			} else {
				expectedBody = string(body)
			}

			assert.Equal(t, tt.requestBody, expectedBody, "request body should be: %s, got: %s",
				tt.requestBody, expectedBody)
			mockHandler.AssertExpectations(t)
		})
	}
}
