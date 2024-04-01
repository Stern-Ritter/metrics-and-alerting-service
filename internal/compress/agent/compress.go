package agent

import (
	"bytes"
	"compress/gzip"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

var compressedContentTypes = []string{"application/json", "text/html"}

func GzipMiddleware(c *resty.Client, resp *resty.Response) error {
	contentType := resp.Header().Values("Content-Type")
	needCompress := utils.Contains(compressedContentTypes, contentType...)

	if needCompress {
		resp.Header().Add("Content-Encoding", "gzip")
		compressedBody, err := compress(resp.Body())
		if err != nil {
			return err
		}
		resp.SetBody(compressedBody)
	}

	return nil
}

func compress(data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	_, err := w.Write(data)
	w.Close()
	return buf.Bytes(), err
}
