package http

import (
	"bytes"
	"encoding/json"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

var b = &backoff.Backoff{
	Min:    2 * time.Second,
	Max:    1 * time.Minute,
	Factor: 2,
	Jitter: true,
}

func Get(path string, output interface{}) (*http.Response, error) {
	for {
		logrus.Info("Request")

		resp, err := http.Get(path)
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}
		}

		if resp != nil && !isServerUnavailable(resp.StatusCode) {
			defer resp.Body.Close()
			if output == nil {
				return resp, nil
			} else {
				return resp, json.NewDecoder(resp.Body).Decode(output)
			}
		}

		time.Sleep(b.Duration())
	}
}

func Post(path, token, jsonStr string) (*http.Response, error) {
	req, err := http.NewRequest("POST", path, bytes.NewBufferString(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	for {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}
		}

		if resp != nil && !isServerUnavailable(resp.StatusCode) {
			defer resp.Body.Close()
			return resp, nil
		}

		time.Sleep(b.Duration())
	}
}

func isServerUnavailable(code int) bool {
	switch code {
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusRequestTimeout:
		return true
	default:
		return false
	}
}
