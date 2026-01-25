package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type msgQueueConnectionMetrics struct {
	errorCounter *prometheus.CounterVec
}

func NewNATSConnectionMetrics(namespace string) ccmetrics.MsgQueueConnectionMetricsHandler {
	return newMsgQueueConnectionMetrics(ccmetrics.NatsSubsystem, namespace)

}

func newMsgQueueConnectionMetrics(subsystem, namespace string) ccmetrics.MsgQueueConnectionMetricsHandler {

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "connection_errors_total",
			Help:      "Total number of message queue connection errors",
		},
		[]string{"type"},
	)

	return &msgQueueConnectionMetrics{
		errorCounter: errorCounter,
	}
}

func (mh *msgQueueConnectionMetrics) ErrorInc(connectionType string) {
	mh.errorCounter.WithLabelValues(connectionType).Inc()
}
