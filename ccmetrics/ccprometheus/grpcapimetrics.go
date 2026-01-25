package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type grpcApiMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestCounter  *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
}

func NewGrpcApiMetricsHandler(namespace string) ccmetrics.GrpcApiMetricsHandler {
	mh := &grpcApiMetrics{}

	mh.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.GrpcInSubsystem,
			Name:      "requests_duration",
			Help:      "gRPC API request latencies",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"path", "code"},
	)

	mh.requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.GrpcInSubsystem,
			Name:      "requests_count",
			Help:      "Total number of gRPC API requests",
		},
		[]string{"path", "code"},
	)

	mh.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: ccmetrics.GrpcInSubsystem,
			Name:      "grpc_api_requests_errors_total",
			Help:      "Total number of gRPC API request errors",
		},
		[]string{"path", "code"},
	)

	return mh
}

func (mh *grpcApiMetrics) RequestDuration(path, code string, duration float64) {
	mh.requestDuration.WithLabelValues(path, code).Observe(duration)
}

func (mh *grpcApiMetrics) RequestCounterInc(path, code string) {
	mh.requestCounter.WithLabelValues(path, code).Inc()
}

func (mh *grpcApiMetrics) ErrorInc(path, code string) {
	mh.errorCounter.WithLabelValues(path, code).Inc()
}
