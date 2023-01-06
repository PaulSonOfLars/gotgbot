package ext

import (
	"fmt"
	"time"
)

// WebhookOpts represent various fields that are needed for configuring the local webhook server.
type WebhookOpts struct {
	// Listen is the address to listen on (eg: localhost, 0.0.0.0, etc).
	Listen string
	// Port is the port listen on (eg 443, 8443, etc).
	Port int
	// ReadTimeout is passed to the http server to limit the time it takes to read an incoming request.
	// See http.Server for more details.
	ReadTimeout time.Duration
	// ReadHeaderTimeout is passed to the http server to limit the time it takes to read the headers of an incoming
	// request.
	// See http.Server for more details.
	ReadHeaderTimeout time.Duration

	// HTTPS cert and key files for custom signed certificates
	CertFile string
	KeyFile  string

	// SecretToken to be used by the bots on this webhook. Used as a security measure to ensure that you set the webhook.
	SecretToken string
}

// GetListenAddr returns the local listening address, including port.
func (w *WebhookOpts) GetListenAddr() string {
	if w.Listen == "" {
		w.Listen = "0.0.0.0"
	}
	if w.Port == 0 {
		w.Port = 443
	}
	return fmt.Sprintf("%s:%d", w.Listen, w.Port)
}
