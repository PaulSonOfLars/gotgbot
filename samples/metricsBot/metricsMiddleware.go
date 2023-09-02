package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	totalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gotgbot",
			Name:      "requests_total",
			Help:      "Number of requests made to the bot API.",
		},
		[]string{
			"api_method",
		},
	)
	totalHTTPErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gotgbot",
			Name:      "http_request_errors_total",
			Help:      "Number of HTTP errors obtained.",
		},
		[]string{
			"api_method",
		},
	)
	totalAPIErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gotgbot",
			Name:      "api_request_errors_total",
			Help:      "Number of bot API errors obtained.",
		},
		[]string{
			"api_method",
			"api_status_code",
			"description",
		},
	)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gotgbot",
			Name:      "api_request_time_seconds",
			Help:      "Duration of requests made to the bot API.",
		},
		[]string{
			"api_method",
		},
	)
)

// Define middleware BotClient.
type metricsBotClient struct {
	// Inline existing client to call, allowing us to chain middlewares.
	// Inlining also avoids us having to redefine helper methods part of the interface.
	gotgbot.BotClient
}

// Define wrapper around existing RequestWithContext method.
// Note: this is the only method that needs redefining.
func (b metricsBotClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.NamedReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	totalRequests.WithLabelValues(method).Inc()
	timer := prometheus.NewTimer(requestDuration.With(prometheus.Labels{
		"http_method": http.MethodPost,
		"api_method":  method,
	}))

	val, err := b.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	timer.ObserveDuration()
	if err != nil {
		tgErr := &gotgbot.TelegramError{}
		if errors.As(err, &tgErr) {
			totalAPIErrors.WithLabelValues(method, strconv.Itoa(tgErr.Code), tgErr.Description).Inc()
		} else {
			totalHTTPErrors.WithLabelValues(method).Inc()
		}
	}
	return val, err
}

func newMetricsClient() metricsBotClient {
	return metricsBotClient{
		BotClient: &gotgbot.BaseBotClient{
			Client:             http.Client{},
			UseTestEnvironment: false,
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	}
}
