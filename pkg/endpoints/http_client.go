package endpoints

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	sharedClient *http.Client
	clientOnce   sync.Once
)

func getHTTPClient() *http.Client {
	clientOnce.Do(func() {
		tr := &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DialContext:         (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
		}
		sharedClient = &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
	})
	return sharedClient
}
