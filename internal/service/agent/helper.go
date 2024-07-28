package agent

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

func sendPostRequest(client *gentleman.Client, endpoint string, headers map[string]string, data []byte) (
	*gentleman.Response, error) {
	req := client.Request()
	req.Method("POST")
	req.Path(endpoint)
	setHeaders(req, headers)
	req.Use(body.JSON(data))

	return req.Send()
}

func setHeaders(req *gentleman.Request, headers map[string]string) {
	for name, value := range headers {
		req.SetHeader(name, value)
	}
}

func getIPAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no ip address found")
}

// SetInterval schedules a task at a specified interval until the context is done.
// It uses a WaitGroup to manage the lifecycle of the task.
func SetInterval(ctx context.Context, wg *sync.WaitGroup, task func(), interval time.Duration) {
	wg.Add(1)

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
