package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"gopkg.in/h2non/gentleman.v2/context"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// compressedContentTypes holds the content types that should be compressed.
var compressedContentTypes = []string{"application/json", "text/html"}

// GzipMiddleware compresses the request body if it contains one of the specified content types.
// It adds the "Content-Encoding: gzip" header to the request and compresses the body.
func GzipMiddleware(ctx *context.Context, h context.Handler) {
	contentType := ctx.Request.Header.Values("Content-Type")
	needCompress := utils.Contains(compressedContentTypes, contentType...)

	if needCompress {
		ctx.Request.Header.Add("Content-Encoding", "gzip")
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error = fmt.Errorf("middleware body compress error: %w", err)
		}
		compressedBody, err := compress(body)
		if err != nil {
			ctx.Error = fmt.Errorf("middleware body compress error: %w", err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewReader(compressedBody))
		ctx.Request.ContentLength = int64(len(compressedBody))
	}
	h.Next(ctx)
}

func compress(data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	_, err := w.Write(data)
	w.Close()
	return buf.Bytes(), err
}
