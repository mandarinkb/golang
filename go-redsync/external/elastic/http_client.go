package elastic

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/mandarinkb/go-redsync/config"
)

var (
	HTTPESClient *http.Client
)

func NewHttpESClient() *http.Client {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.C().Elastic.HTTPTransport.SkipVerifyTLS},
		Dial: (&net.Dialer{
			Timeout:   config.C().Elastic.HTTPTransport.DialTimeout,
			KeepAlive: config.C().Elastic.HTTPTransport.DialKeepAlive,
		}).Dial,
		MaxIdleConns:          config.C().Elastic.HTTPTransport.MaxIdleConns,
		MaxIdleConnsPerHost:   config.C().Elastic.HTTPTransport.MaxIdleConnsPerHost,
		IdleConnTimeout:       config.C().Elastic.HTTPTransport.IdleConnTimeout,
		TLSHandshakeTimeout:   config.C().Elastic.HTTPTransport.TLSHandshakeTimeout,
		ResponseHeaderTimeout: config.C().Elastic.HTTPTransport.ResponseHeaderTimeout,
		ExpectContinueTimeout: config.C().Elastic.HTTPTransport.ExpectContinueTimeout,
	}
	return &http.Client{
		Transport: httpTransport,
		Timeout:   config.C().Elastic.HTTPTransport.Timeout,
	}
}
