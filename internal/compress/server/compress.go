package server

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// compressedContentTypes holds the content types that should be compressed.
var compressedContentTypes = []string{"application/json", "text/html"}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter is constructor for creating a new compressWriter wrapping http.ResponseWriter
// for writing compressed response body.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns response headers.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the data to response body.
func (c *compressWriter) Write(body []byte) (int, error) {
	contentType := c.Header().Values("Content-type")
	needCompress := utils.Contains(contentType, compressedContentTypes...)
	if needCompress {
		c.w.Header().Set("Content-Encoding", "gzip")
		return c.zw.Write(body)
	}
	return c.w.Write(body)
}

// WriteHeader write response headers and set response status code.
func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip.Writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader is constructor for creating a new compressReader for reading compressed request body
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads data into the provided byte slice and returns the number of bytes read.
func (c *compressReader) Read(body []byte) (n int, err error) {
	return c.zr.Read(body)
}

// Close closes the gzip.Reader and the io.ReadCloser.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware is an HTTP middleware for reading compressed request body and
// writing compressed response body if the client supports gzip encoding.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Values("Content-Encoding")
		sendsGzip := utils.Contains(contentEncoding, "gzip")
		contentType := r.Header.Values("Content-type")
		needUncompressed := utils.Contains(compressedContentTypes, contentType...)

		if sendsGzip && needUncompressed {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			r.Body = cr
		}

		acceptEncoding := r.Header.Values("Accept-Encoding")
		supportsGzip := utils.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			cw := NewCompressWriter(w)
			defer cw.Close()
			w = cw
		}

		next.ServeHTTP(w, r)
	})
}
