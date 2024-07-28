package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
)

const (
	signKey = "HashSHA256"
)

// signWriter wraps http.ResponseWriter and adds an HMAC SHA256 signature to the response body.
type signWriter struct {
	w         http.ResponseWriter
	secretKey string
}

// NewSignWriter is constructor for creating a new signWriter.
func NewSignWriter(w http.ResponseWriter, secretKey string) *signWriter {
	return &signWriter{w: w, secretKey: secretKey}
}

// Header returns response headers.
func (s *signWriter) Header() http.Header {
	return s.w.Header()
}

// Write writes the data to response body.
func (s *signWriter) Write(body []byte) (int, error) {
	sign := getSign(body, s.secretKey)
	s.w.Header().Set(signKey, sign)
	return s.w.Write(body)
}

// WriteHeader write response headers and set response status code.
func (s *signWriter) WriteHeader(statusCode int) {
	s.w.WriteHeader(statusCode)
}

// SignMiddleware is a middleware that checks the request body signature and
// signs the response body with HMAC SHA256 if a secret key is configured.
func (s *Server) SignMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sign := r.Header.Get(signKey)

		hasSign := len(strings.TrimSpace(sign)) > 0
		needCheckSign := len(s.Config.SecretKey) > 0
		needSignResponseBody := len(s.Config.SecretKey) > 0

		if hasSign && needCheckSign {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Read request body error", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			if len(body) > 0 {
				err = checkSign(body, sign, s.Config.SecretKey)
				if err != nil {
					http.Error(w, "Invalid request body sign", http.StatusBadRequest)
					return
				}
			}
		}

		ow := w
		if needSignResponseBody {
			ow = NewSignWriter(w, s.Config.SecretKey)
		}

		next.ServeHTTP(ow, r)
	})
}

func checkSign(value []byte, sign string, secretKey string) error {
	decodedSign, err := hex.DecodeString(sign)
	if err != nil {
		return errors.NewUnsignedRequest("Invalid request body sign", nil)
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)

	if !hmac.Equal(decodedSign, hash) {
		return errors.NewUnsignedRequest("Invalid request body sign", nil)
	}

	return nil
}

func getSign(value []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}
