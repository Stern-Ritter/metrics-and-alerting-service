package agent

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

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

func sendPostRequest(client *resty.Client, url, endpoint, contentType string, pathParams map[string]string) (*resty.Response, error) {
	resp, err := client.R().
		SetHeader("Content-Type", contentType).
		SetPathParams(pathParams).
		Get(utils.AddProtocolPrefix(strings.Join([]string{url, endpoint}, "")))

	return resp, err
}
