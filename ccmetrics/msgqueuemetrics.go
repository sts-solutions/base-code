package ccmetrics

const (
	NatsSubsystem = "nats"
)

type MsgQueueConnectionMetricsHandler interface {
	ErrorInc(errorType string)
}

type MsgQueueConsumerMetricsHandler interface {
	ConcurrentProcessInc(subject, consumerName string, counter int)
	MessageDuration(subject, consumerName string, duration float64)
	MessageCounterInc(subject, consumerName string)
	ErrorInc(errorType string)
}

type MsgQueuePublisherMetricsHandler interface {
	MessageDuration(subject string, duration float64)
	MessageCounterInc(subject string)
	ErrorInc(errorType string)
}
