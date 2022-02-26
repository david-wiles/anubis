package anubis

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Option interface {
	SetOpt(*Anubis)
}

type OutputOpt string

func (opt OutputOpt) SetOpt(anubis *Anubis) { anubis.Output = string(opt) }

type HeaderOpt struct {
	Key, Value string
}

func (opt HeaderOpt) SetOpt(anubis *Anubis) {
	anubis.Headers[opt.Key] = opt.Value
}

type NWorkerOpt int

func (opt NWorkerOpt) SetOpt(anubis *Anubis) { anubis.Workers = int(opt) }

// ProxyOpt allows for setting the network proxy for the system. The proxy URL should be passed as a string
// and will be used to construct a new client using the DefaultWebDriver.
// If a different web driver should be used, then the proxy should be configured using the WebDriverOpt instead.
type ProxyOpt string

func (opt ProxyOpt) SetOpt(anubis *Anubis) {
	// If proxy is an empty string, this is a nop
	if opt == "" {
		return
	}

	proxy := http.ProxyFromEnvironment
	u, err := url.Parse(string(opt))
	if err != nil {
		log.Println(err, "Using proxy from environment instead")
	}

	proxy = http.ProxyURL(u)

	driver := DefaultWebDriver{
		client: http.Client{
			Transport: &http.Transport{
				Proxy: proxy,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}

	anubis.Driver = driver
}

type WebDriverOpt struct {
	Driver WebDriver
}

func (opt WebDriverOpt) SetOpt(anubis *Anubis) { anubis.Driver = opt.Driver }

type ResponseHandlerOpt struct {
	Handler ResponseHandler
}

func (opt ResponseHandlerOpt) SetOpt(anubis *Anubis) { anubis.Handler = opt.Handler }
