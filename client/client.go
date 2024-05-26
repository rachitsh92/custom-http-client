package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Config struct {
	Timeout int
	Retries int
}

type HTTPClient struct {
	Client  *http.Client
	Retries int
}

func NewClient(config Config) *HTTPClient {
	return &HTTPClient{

		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13,
				},
				TLSHandshakeTimeout: 10*time.Second,
			},
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		Retries: config.Retries,
	}
}
