package ccmetrics

import (
	"github.com/sts-solutions/base-code/ccmetrics"

	"github.com/prometheus/client_golang/prometheus"
)

type databaseMetricHandler struct {
	requestDuration           *prometheus.HistogramVec
	requestCounter            *prometheus.CounterVec
	errorCounter              *prometheus.CounterVec
	maxOpenConnectionsCounter *prometheus.GaugeVec
	openConnectionsCounter    *prometheus.GaugeVec
	waitCounter               *prometheus.GaugeVec
	waitDuration              *prometheus.GaugeVec
	idleCounter               *prometheus.GaugeVec
	inUseCounter              *prometheus.GaugeVec
	queueCounter              *prometheus.GaugeVec
}

func NewDatabaseMetricHandler(namespace string, subsystem string) ccmetrics.DatabaseMetricsHandler {
	mh := &databaseMetricHandler{}

	mh.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "queries_duration",
			Help:      "Query latencies",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"source"},
	)

	mh.requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "queries_total",
			Help:      "Total number of queries",
		},
		[]string{"source"},
	)

	mh.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "queries_errors_total",
			Help:      "Total number of query errors",
		},
		[]string{"source"},
	)

	mh.maxOpenConnectionsCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "max_open_connections",
			Help:      "Maximum number of open connections to the database",
		},
		[]string{},
	)

	mh.openConnectionsCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "open_connections",
			Help:      "Number of open connections to the database",
		},
		[]string{},
	)

	mh.waitCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "wait_count",
			Help:      "Number of waits for a connection from the pool",
		},
		[]string{},
	)

	mh.waitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "wait_duration_seconds",
			Help:      "Total time waited for a connection from the pool",
		},
		[]string{},
	)

	mh.idleCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "idle_connections",
			Help:      "Number of idle connections in the pool",
		},
		[]string{"source"},
	)

	mh.inUseCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "in_use_connections",
			Help:      "Number of in-use connections in the pool",
		},
		[]string{},
	)

	mh.queueCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "connection_queue_length",
			Help:      "Number of waiting requests for a connection from the pool",
		},
		[]string{"name"},
	)

	return mh
}

func (m *databaseMetricHandler) RequestDuration(source string, duration float64) {
	m.requestDuration.WithLabelValues(source).Observe(duration)
}

func (m *databaseMetricHandler) RequestCounterInc(source string) {
	m.requestCounter.WithLabelValues(source).Inc()
}

func (m *databaseMetricHandler) ErrorInc(source string) {
	m.errorCounter.WithLabelValues(source).Inc()
}

func (m *databaseMetricHandler) UpdateMaxOpenConnections(count int) {
	m.maxOpenConnectionsCounter.WithLabelValues().Set(float64(count))
}

func (m *databaseMetricHandler) UpdateOpenConnections(count int) {
	m.openConnectionsCounter.WithLabelValues().Set(float64(count))
}

func (m *databaseMetricHandler) UpdateWaitDuration(totalDuration int) {
	m.waitDuration.WithLabelValues().Set(float64(totalDuration))
}

func (m *databaseMetricHandler) UpdateWaitCount(count int) {
	m.waitCounter.WithLabelValues().Set(float64(count))
}

func (m *databaseMetricHandler) UpdateIdleCount(count int) {
	m.idleCounter.WithLabelValues().Set(float64(count))
}

func (m *databaseMetricHandler) UpdateInUseCount(count int) {
	m.inUseCounter.WithLabelValues().Set(float64(count))
}

func (m *databaseMetricHandler) QueueCount(name string, count int) {
	m.queueCounter.WithLabelValues(name).Set(float64(count))
}
