package ext

import (
	"time"
)

// WebhookOpts represent various fields that are needed for configuring the local webhook server.
type WebhookOpts struct {
	// ListenAddr is the address and port to listen on (eg: localhost:http, 0.0.0.0:8080, :https, "[::1]:", etc).
	// See the net package for details.
	ListenAddr string
	// ListenNet is the network type to listen on (must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket").
	// Empty means the default, "tcp".
	ListenNet string
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

func (w *WebhookOpts) GetListenNet() string {
	if w.ListenNet == "" {
		return "tcp"
	}
	return w.ListenNet
}
