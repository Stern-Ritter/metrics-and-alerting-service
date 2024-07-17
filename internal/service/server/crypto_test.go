package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptMiddleware(t *testing.T) {
	type request struct {
		isBodyEncrypted bool
		body            string
	}

	type response struct {
		status int
		body   string
	}

	tests := []struct {
		name          string
		hasPrivateKey bool
		request       request
		response      response
	}{
		{
			name:          "should decrypt request body when private key is defined and request body isn`t empty",
			hasPrivateKey: true,
			request: request{
				isBodyEncrypted: true,
				body:            "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusOK,
				body:   "The Ultimate Question of Life, the Universe, and Everything",
			},
		},
		{
			name:          "shouldn`t decrypt request body when private key is defined and request body is empty",
			hasPrivateKey: true,
			request: request{
				isBodyEncrypted: false,
				body:            "",
			},
			response: response{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name:          "shouldn`t decrypt request body when private key isn`t defined and request body isn`t empty",
			hasPrivateKey: false,
			request: request{
				isBodyEncrypted: false,
				body:            "The Ultimate Question of Life, the Universe, and Everything",
			},
			response: response{
				status: http.StatusOK,
				body:   "The Ultimate Question of Life, the Universe, and Everything",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err, "unexpected error when generate private rsa key")

			var rsaPrivateKey *rsa.PrivateKey
			if tt.hasPrivateKey {
				rsaPrivateKey = privateKey
			}

			var requestBody []byte
			if tt.request.isBodyEncrypted {
				requestBody, err = rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte(tt.request.body))
				require.NoError(t, err, "unexpected error when encrypt request body")
			} else {
				requestBody = []byte(tt.request.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestBody)))

			server := &Server{
				rsaPrivateKey: rsaPrivateKey,
			}

			handler := server.EncryptMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				_, _ = w.Write(body)
			}))

			r := httptest.NewRecorder()
			handler.ServeHTTP(r, req)

			resp := r.Result()
			data, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "unexpected error when read response body")

			gotStatus := resp.StatusCode
			gotBody := string(data)
			assert.Equal(t, tt.response.status, gotStatus, "response status code should be %d, got %d",
				tt.response.status, gotStatus)
			assert.Equal(t, tt.response.body, gotBody, "response body should be: %s, got: %s",
				tt.response.body, gotBody)
		})
	}
}
