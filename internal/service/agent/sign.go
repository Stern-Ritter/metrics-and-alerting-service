package agent

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
)

const (
	signKey = "HashSHA256"
)

// SignMiddleware is a middleware that signs the request body with HMAC SHA256 if a secret key is configured.
func (a *Agent) SignMiddleware(ctx *gcontext.Context, h gcontext.Handler) {
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

// SignInterceptor is a gRPC client interceptor that signs the request with a secret key.
func (a *Agent) SignInterceptor(ctx context.Context, method string, req interface{}, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	needSignRequest := len(strings.TrimSpace(a.Config.SecretKey)) > 0
	if needSignRequest {
		message, ok := req.(proto.Message)
		if !ok {
			return status.Errorf(codes.Internal, "sign interceptor: request isn't a proto.Message")
		}

		body, err := proto.Marshal(message)
		if err != nil {
			return status.Errorf(codes.Internal, "sign interceptor: %s", err)
		}

		if len(body) > 0 {
			sign := getSign(body, a.Config.SecretKey)
			ctx = metadata.AppendToOutgoingContext(ctx, signKey, sign)
		}
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func getSign(value []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(value)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
