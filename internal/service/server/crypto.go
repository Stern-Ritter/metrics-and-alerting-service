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
		isEncryptionEnabled := s.rsaPrivateKey != nil
		if isEncryptionEnabled {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Read request body error", http.StatusBadRequest)
				return
			}

			if len(body) > 0 {
				decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, s.rsaPrivateKey, body)
				if err != nil {
					http.Error(w, "Decrypt request body error", http.StatusBadRequest)
				}

				r.Body = io.NopCloser(bytes.NewReader(decryptedBody))
				r.ContentLength = int64(len(decryptedBody))
			}
		}

		next.ServeHTTP(w, r)
	})
}
