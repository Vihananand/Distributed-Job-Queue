package utils

import (
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func InitLogger(env string) {
	var h slog.Handler
	if env == "production" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(h))
}

type Metrics struct {
	JobsPushed    prometheus.Counter
	JobsProcessed prometheus.Counter
	JobsFailed    prometheus.Counter
	QueueLength   prometheus.Gauge
	Registry *prometheus.Registry
}

func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()

	pushed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "djq_jobs_pushed_total",
		Help: "Total number of jobs submitted to the queue.",
	})
	processed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "djq_jobs_processed_total",
		Help: "Total number of jobs successfully processed.",
	})
	failed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "djq_jobs_failed_total",
		Help: "Total number of job execution failures (including retried).",
	})
	length := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "djq_queue_length",
		Help: "Current number of pending (not-yet-running) jobs.",
	})

	reg.MustRegister(pushed, processed, failed, length)

	return &Metrics{
		JobsPushed:    pushed,
		JobsProcessed: processed,
		JobsFailed:    failed,
		QueueLength:   length,
		Registry:      reg,
	}
}
