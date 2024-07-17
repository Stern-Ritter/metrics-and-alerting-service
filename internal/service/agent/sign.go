package agent

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"gopkg.in/h2non/gentleman.v2/context"
)

const (
	signKey = "HashSHA256"
)

// SignMiddleware is a middleware that signs the request body with HMAC SHA256 if a secret key is configured.
func (a *Agent) SignMiddleware(ctx *context.Context, h context.Handler) {
	needSignResponseBody := len(strings.TrimSpace(a.Config.SecretKey)) > 0
	if needSignResponseBody {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error = fmt.Errorf("middleware body sign error: %w", err)
			h.Next(ctx)
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

		if len(body) > 0 {
			sign := getSign(body, a.Config.SecretKey)
			ctx.Request.Header.Add(signKey, sign)
		}
	}

	h.Next(ctx)
}

func getSign(value []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
