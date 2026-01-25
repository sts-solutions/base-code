package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type publisherMetrics struct {
	messageDuration *prometheus.HistogramVec
	messageCounter  *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
}

func NewNATSMsgQueuePublishMetricsHandler(namespace string) ccmetrics.MsgQueuePublisherMetricsHandler {
	return newMsgQueuePublishMetricsHandler(ccmetrics.NatsSubsystem, namespace)
}

func newMsgQueuePublishMetricsHandler(subsystem, namespace string) ccmetrics.MsgQueuePublisherMetricsHandler {

	messageDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published_message_duration",
			Help:      "Message processing durations in seconds",
			Buckets:   ccmetrics.DefaultBuckets,
		},
		[]string{"subject", "name"},
	)

	messageCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published_messages_count",
			Help:      "Total number of messages processed",
		},
		[]string{"subject", "name"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "published§_messages_errors_total",
			Help:      "Total number of message processing errors",
		},
		[]string{"subject", "name"},
	)

	return &publisherMetrics{
		messageDuration: messageDuration,
		messageCounter:  messageCounter,
		errorCounter:    errorCounter,
	}
}

func (mh *publisherMetrics) MessageDuration(subject string, duration float64) {
	mh.messageDuration.WithLabelValues(subject).Observe(duration)
}

func (mh *publisherMetrics) MessageCounterInc(subject string) {
	mh.messageCounter.WithLabelValues(subject).Inc()
}

func (mh *publisherMetrics) ErrorInc(errorType string) {
	mh.errorCounter.WithLabelValues(errorType).Inc()
}
