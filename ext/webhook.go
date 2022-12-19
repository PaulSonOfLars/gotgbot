package ext

import (
	"fmt"
	"strings"
	"time"
)

// WebhookOpts represent various fields that are needed for configuring the local webhook server.
type WebhookOpts struct {
	// Listen is the address to listen on (eg: localhost, 0.0.0.0, etc).
	Listen string
	// Port is the port listen on (eg 443, 8443, etc).
	Port int
	// URLPath defines the path to listen at; eg <domainname>/<URLPath>.
	// Using the bot token here is often a good idea, as it is a secret known only by telegram.
	URLPath string
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

	// The secret token used in the Bot.SetWebhook call, which can be used to ensure that the request comes from a
	// webhook set by you.
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

// GetWebhookURL returns the domain in the form domain/path.
// eg: example.com/super_secret_token
func (w *WebhookOpts) GetWebhookURL(domain string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(domain, "/"), w.URLPath)
}
