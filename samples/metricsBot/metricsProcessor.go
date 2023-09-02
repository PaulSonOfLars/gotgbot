package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	totalUpdates = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "gotgbot",
			Name:      "updates_total",
			Help:      "Number of incoming updates.",
		},
	)
	updateProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "gotgbot",
			Name:      "update_processing_time_seconds",
			Help:      "Time to process each update.",
		},
	)
	bufferedUpdates = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gotgbot",
			Name:      "buffered_updates",
			Help:      "Number of updates currently buffered in the dispatcher limiter channel.",
		},
	)
	bufferedUpdatesLimit = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gotgbot",
			Name:      "buffered_updates_limit",
			Help:      "Maximum number of buffered updates in the limiter channel.",
		},
	)
)

var _ ext.Processor = metricsProcessor{}

type metricsProcessor struct {
	processor ext.Processor
}

func (m metricsProcessor) ProcessUpdate(d *ext.Dispatcher, b *gotgbot.Bot, ctx *ext.Context) error {
	totalUpdates.Inc()
	timer := prometheus.NewTimer(updateProcessingDuration)
	defer timer.ObserveDuration()

	return m.processor.ProcessUpdate(d, b, ctx)
}

func monitorDispatcherBuffer(d *ext.Dispatcher) {
	// Limiter cant be changed during execution, so only needs to be set once.
	bufferedUpdatesLimit.Set(float64(d.MaxUsage()))

	for {
		bufferedUpdates.Set(float64(d.CurrentUsage()))
		time.Sleep(time.Second)
	}
}
