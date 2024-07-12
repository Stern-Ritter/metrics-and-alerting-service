package agent

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gentleman.v2/context"
)

type mockHandler struct{}

func (m *mockHandler) Next(ctx *context.Context)             {}
func (m *mockHandler) Stop(ctx *context.Context)             {}
func (m *mockHandler) Error(ctx *context.Context, err error) {}

func TestGzipMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		contentType      string
		body             string
		expectCompressed bool
		expectError      bool
	}{
		{
			name:             "should compress JSON",
			contentType:      "application/json",
			body:             `{"answer": 42}`,
			expectCompressed: true,
		},
		{
			name:             "should compress HTML",
			contentType:      "text/html",
			body:             `<html><body>The Ultimate Question of Life, the Universe and Everything.</body></html>`,
			expectCompressed: true,
		},
		{
			name:             "should not compress plain text",
			contentType:      "text/plain",
			body:             "The Ultimate Question of Life, the Universe and Everything.",
			expectCompressed: false,
		},
		{
			name:             "should compress empty body with compressible type",
			contentType:      "application/json",
			body:             ``,
			expectCompressed: true,
		},
		{
			name:             "should not compress empty body with not compressible type",
			contentType:      "text/plain",
			body:             ``,
			expectCompressed: false,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.New()
			ctx.Request.Header.Set("Content-Type", tt.contentType)
			ctx.Request.Body = io.NopCloser(strings.NewReader(tt.body))

			GzipMiddleware(ctx, &mockHandler{})

			if tt.expectError {
				require.Error(t, ctx.Error, "expected error but got none")
			} else {
				require.NoError(t, ctx.Error, "unexpected error: %v", ctx.Error)
			}

			if tt.expectCompressed {
				expectedEncoding := "gzip"
				gotEncoding := ctx.Request.Header.Get("Content-Encoding")
				assert.Equal(t, expectedEncoding, gotEncoding,
					"expected Content-Encoding to be %s, got %s", expectedEncoding, gotEncoding)

				compressedBody, err := io.ReadAll(ctx.Request.Body)
				require.NoError(t, err, "failed to read compressed body: %v", err)
				decompressedBody, err := decompress(compressedBody)
				require.NoError(t, err, "failed to decompress body: %v", err)
				assert.Equal(t, tt.body, string(decompressedBody),
					"expected decompressed body to be %s, got %s", tt.body, decompressedBody)
			} else {
				expectedEncoding := ""
				gotEncoding := ctx.Request.Header.Get("Content-Encoding")
				assert.Equal(t, expectedEncoding, gotEncoding,
					"expected no Content-Encoding header, got %s", gotEncoding)

				body, err := io.ReadAll(ctx.Request.Body)
				require.NoError(t, err, "failed to read body: %v", err)
				assert.Equal(t, tt.body, string(body), "expected body to be %s, got %s", tt.body, body)
			}
		})
	}
}

func decompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
