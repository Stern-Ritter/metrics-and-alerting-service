package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
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
				body:   "Invalid request sign\n",
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
				_, _ = w.Write([]byte(tt.request.body))
			}))

			r := httptest.NewRecorder()
			handler.ServeHTTP(r, req)

			data, err := io.ReadAll(r.Body)
			require.NoError(t, err, "unexpected error reading body")

			gotStatus := r.Code
			gotSign := r.Header().Get(signKey)
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
