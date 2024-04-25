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

type signWriter struct {
	w         http.ResponseWriter
	secretKey string
}

func NewSignWriter(w http.ResponseWriter, secretKey string) *signWriter {
	return &signWriter{w: w, secretKey: secretKey}
}

func (s *signWriter) Header() http.Header {
	return s.w.Header()
}

func (s *signWriter) Write(body []byte) (int, error) {
	sign := getSign(body, s.secretKey)
	s.w.Header().Set("HashSHA256", sign)
	return s.w.Write(body)
}

func (s *signWriter) WriteHeader(statusCode int) {
	s.w.WriteHeader(statusCode)
}

func (s *Server) SignMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sign := r.Header.Get("HashSHA256")
		hasBody := r.Body != http.NoBody
		needCheckSign := len(strings.TrimSpace(s.Config.SecretKey)) != 0 && len(strings.TrimSpace(sign)) != 0

		if hasBody && needCheckSign {
			if len(sign) == 0 {
				http.Error(w, "Unsigned request body", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Read request body error", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			err = checkSign(body, sign, s.Config.SecretKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ow := w
		if needCheckSign {
			ow = NewSignWriter(w, s.Config.SecretKey)
		}

		next.ServeHTTP(ow, r)
	})
}

func checkSign(value []byte, sign string, secretKey string) error {
	decodedSign, err := hex.DecodeString(sign)
	if err != nil {
		return errors.NewUnsignedRequest("Invalid request sign", err)
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)

	if !hmac.Equal(decodedSign, hash) {
		return errors.NewUnsignedRequest("Invalid request sign", err)
	}

	return nil
}

func getSign(value []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}
