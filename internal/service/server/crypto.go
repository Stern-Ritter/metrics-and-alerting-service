package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
)

// EncryptMiddleware is a middleware that decrypts the request body using RSA encryption.
//
// If the RSA private key is set in the Server, this middleware reads the encrypted request body,
// decrypts it using the private key, and replaces the original request body with the decrypted body.
// The decrypted body is then passed along the middleware chain.
func (s *Server) EncryptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasBody := r.Body != http.NoBody
		isEncryptionEnabled := s.rsaPrivateKey != nil

		if hasBody && isEncryptionEnabled {
			encryptedBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Read encrypted request body error", http.StatusBadRequest)
				return
			}
			body, err := rsa.DecryptPKCS1v15(rand.Reader, s.rsaPrivateKey, encryptedBody)
			if err != nil {
				http.Error(w, "Decrypt request body error", http.StatusBadRequest)
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
			r.ContentLength = int64(len(body))
		}

		next.ServeHTTP(w, r)
	})
}
