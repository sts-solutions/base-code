package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type httpApiMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestCounter  *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
}

func NewHttpApiMetricsHandler(namespace string) ccmetrics.HttpApiMetricsHandler {
	mh := &httpApiMetrics{}

	mh.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.HttpInSubsystem,
			Name:      "requests_duration",
			Help:      "HTTP API request latencies",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"method", "path", "code"},
	)

	mh.requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.HttpInSubsystem,
			Name:      "requests_count",
			Help:      "Total number of HTTP API requests",
		},
		[]string{"method", "path", "code"},
	)

	mh.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.HttpInSubsystem,
			Name:      "http_api_requests_errors_total",
			Help:      "Total number of HTTP API request errors",
		},
		[]string{"method", "path", "code"},
	)

	return mh
}

func (mh *httpApiMetrics) RequestDuration(method, path, code string, duration float64) {
	mh.requestDuration.WithLabelValues(method, path, code).Observe(duration)
}

func (mh *httpApiMetrics) RequestCounterInc(method, path, code string) {
	mh.requestCounter.WithLabelValues(method, path, code).Inc()
}

func (mh *httpApiMetrics) ErrorInc(method, path, code string) {
	mh.errorCounter.WithLabelValues(method, path, code).Inc()
}
