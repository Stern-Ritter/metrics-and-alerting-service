package agent

import (
	"context"
	"sync"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

func sendPostRequest(client *gentleman.Client, endpoint, contentType string, data []byte) (*gentleman.Response, error) {
	req := client.Request()
	req.Method("POST")
	req.Path(endpoint)
	req.SetHeader("Content-Type", contentType)
	req.Use(body.JSON(data))

	return req.Send()
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
