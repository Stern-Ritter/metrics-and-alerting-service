package agent

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

func sendPostRequest(client *resty.Client, url, endpoint, contentType string, body []byte) (*resty.Response, error) {
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{contentType}

	resp, err := client.R().
		SetHeaderMultiValues(headers).
		SetBody(body).
		Post(utils.AddProtocolPrefix(strings.Join([]string{url, endpoint}, "")))

	return resp, err
}

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
