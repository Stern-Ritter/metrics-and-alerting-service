package agent

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"

	"gopkg.in/h2non/gentleman.v2/context"

	crypto "github.com/Stern-Ritter/metrics-and-alerting-service/internal/crypto/agent"
)

// GetRSAPublicKey reads an RSA public key from a PEM-encoded file.
//
// This function checks if the provided path is not empty. If the path is
// provided, it attempts to read the RSA public key from the specified file.
func GetRSAPublicKey(path string) (*rsa.PublicKey, error) {
	var rsaPublicKey *rsa.PublicKey

	isEncryptionEnabled := len(path) != 0
	if isEncryptionEnabled {
		key, err := crypto.GetRSAPublicKey(path)
		if err != nil {
			return nil, err
		}
		rsaPublicKey = key
	}

	return rsaPublicKey, nil
}

// EncryptMiddleware is a middleware that encrypts the request body using RSA encryption.
//
// If the RSA public key is set in the Agent, this middleware reads the request body, encrypts it
// using the public key, and replaces the original request body with the encrypted body. The
// encrypted body is then passed along the middleware chain.
func (a *Agent) EncryptMiddleware(ctx *context.Context, h context.Handler) {
	isEncryptionEnabled := a.rsaPublicKey != nil
	if isEncryptionEnabled {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error = fmt.Errorf("middleware body encryption error: %w", err)
			h.Next(ctx)
			return
		}

		if len(body) > 0 {
			encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, a.rsaPublicKey, body)
			if err != nil {
				ctx.Error = fmt.Errorf("middleware body encryption error: %w", err)
				h.Next(ctx)
				return
			}

			ctx.Request.Body = io.NopCloser(bytes.NewReader(encryptedBody))
			ctx.Request.ContentLength = int64(len(encryptedBody))
		}
	}

	h.Next(ctx)
}
