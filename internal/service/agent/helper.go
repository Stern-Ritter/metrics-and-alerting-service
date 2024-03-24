package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

var compressedContentTypes = []string{"application/json", "text/html"}

func SetInterval(ctx context.Context, wg *sync.WaitGroup, task func(), interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				task()
				time.Sleep(interval)
			}
		}
	}()
}

func sendPostRequest(client *resty.Client, url, endpoint, contentType string, body []byte) (*resty.Response, error) {
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{contentType}

	sendedBody := body
	needCompress := utils.Contains(compressedContentTypes, contentType)
	if needCompress {
		headers["Content-Encoding"] = []string{"gzip"}
		compressedBody, err := compress(body)
		if err != nil {
			return nil, err
		}
		sendedBody = compressedBody
	}

	resp, err := client.R().
		SetHeaderMultiValues(headers).
		SetBody(sendedBody).
		Post(utils.AddProtocolPrefix(strings.Join([]string{url, endpoint}, "")))

	return resp, err
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	w.Close()
	return b.Bytes(), err
}
