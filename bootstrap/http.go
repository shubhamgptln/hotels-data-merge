package bootstrap

import (
	"net"
	"net/http"
	"time"
)

func NewHTTPClient(c *Config) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          10,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: c.HTTP.Timeout,
			DialContext: (&net.Dialer{
				Timeout:   c.HTTP.Timeout,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}

	return client
}
