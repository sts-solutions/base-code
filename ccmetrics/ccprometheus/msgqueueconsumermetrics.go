package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type consumerMetrics struct {
	concurrencyProcess *prometheus.GaugeVec
	messageDuration    *prometheus.HistogramVec
	messageCounter     *prometheus.CounterVec
	errorCounter       *prometheus.CounterVec
}

func NewNATSMsgQueueConsumeMetricsHandler(namespace string) ccmetrics.MsgQueueConsumerMetricsHandler {
	return newMsgQueueConsumerMetricsHandler(ccmetrics.NatsSubsystem, namespace)
}

func newMsgQueueConsumerMetricsHandler(subsystem, namespace string) ccmetrics.MsgQueueConsumerMetricsHandler {

	concurrencyProcess := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "concurrency_process",
			Help:      "Number of concurrent message processing",
		},
		[]string{"subject", "name"},
	)

	messageDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "consumed_message_duration",
			Help:      "Message processing durations in seconds",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"subject", "name"},
	)

	messageCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "consumed_messages_count",
			Help:      "Total number of messages processed",
		},
		[]string{"subject", "name"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "consumed_messages_errors_total",
			Help:      "Total number of message processing errors",
		},
		[]string{"subject", "name"},
	)

	return &consumerMetrics{
		concurrencyProcess: concurrencyProcess,
		messageDuration:    messageDuration,
		messageCounter:     messageCounter,
		errorCounter:       errorCounter,
	}
}

func (mh *consumerMetrics) ConcurrentProcessInc(subject, name string, quantity int) {
	mh.concurrencyProcess.WithLabelValues(subject, name).Set(float64(quantity))
}

func (mh *consumerMetrics) MessageDuration(subject, name string, duration float64) {
	mh.messageDuration.WithLabelValues(subject, name).Observe(duration)
}

func (mh *consumerMetrics) MessageCounterInc(subject, name string) {
	mh.messageCounter.WithLabelValues(subject, name).Inc()
}

func (mh *consumerMetrics) ErrorInc(errorType string) {
	mh.errorCounter.WithLabelValues(errorType).Inc()
}
