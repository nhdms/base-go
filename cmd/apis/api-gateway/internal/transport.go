package internal

import (
	"github.com/nhdms/base-go/pkg/config"
	"net"
	"net/http"
	"time"
)

func getDefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   config.ViperGetDurationWithDefault("transport.dial_timeout", 5*time.Second), // Reduced from 30s
			KeepAlive: config.ViperGetDurationWithDefault("transport.dial_keepalive", 30*time.Second),
		}).DialContext,
		TLSHandshakeTimeout:   config.ViperGetDurationWithDefault("transport.tls_handshake", 5*time.Second),    // Reduced from 10s
		ResponseHeaderTimeout: config.ViperGetDurationWithDefault("transport.response_header", 10*time.Second), // Reduced from 30s
		IdleConnTimeout:       config.ViperGetDurationWithDefault("transport.idle_timeout", 90*time.Second),
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          config.ViperGetIntWithDefault("transport.max_idle_conns", 1000), // Increased from 500
		MaxIdleConnsPerHost:   config.ViperGetIntWithDefault("transport.max_idle_conns_per_host", MaxIdleConnectPerHost),
		ForceAttemptHTTP2:     true, // Enable HTTP/2
		DisableCompression:    true, // Disable compression handling
	}
}
