package ccmetrics

import (
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
	"google.golang.org/grpc"
)

type downstreamGRPCMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestCounter  *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
	interceptor     grpc.UnaryClientInterceptor
}

func NewDownstreamGRPCMetrics(namespace string) *downstreamGRPCMetrics {
	mh := &downstreamGRPCMetrics{}

	mh.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "requests_duration",
			Help:      "Downstream gRPC request latencies",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"recipient", "method", "code"},
	)

	mh.requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_count",
			Help:      "Total number of downstream gRPC requests",
		},
		[]string{"recipient", "method", "code"},
	)

	mh.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "errors_count",
			Help:      "Total number of downstream gRPC request errors",
		},
		[]string{"recipient", "method", "code"},
	)

	grpcMetric := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(ccmetrics.GrpcClientDefaultBuckets),
		),
	)

	prometheus.DefaultRegisterer.MustRegister(grpcMetric)
	mh.interceptor = grpcMetric.UnaryClientInterceptor()

	return mh
}

func (mh *downstreamGRPCMetrics) RequestDuration(recipient, method, code string, duration float64) {
	mh.requestDuration.WithLabelValues(recipient, method, code).Observe(duration)
}

func (mh *downstreamGRPCMetrics) RequestCounterInc(recipient, method, code string) {
	mh.requestCounter.WithLabelValues(recipient, method, code).Inc()
}

func (mh *downstreamGRPCMetrics) ErrorInc(recipient, method, code string) {
	mh.errorCounter.WithLabelValues(recipient, method, code).Inc()
}

func (mh *downstreamGRPCMetrics) Interceptor() grpc.UnaryClientInterceptor {
	return mh.interceptor
}
