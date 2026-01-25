package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type counter struct {
	counter *prometheus.GaugeVec
}

func NewCounter(namespace string, name string, help string, labelNames []string) ccmetrics.Counter {
	c := &counter{}
	c.counter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		},
		labelNames,
	)
	return c
}
func (c *counter) Inc(labels ...string) {
	c.counter.WithLabelValues(labels...).Inc()
}

func (c *counter) Set(value float64, labels ...string) {
	c.counter.WithLabelValues(labels...).Set(value)
}

func (c *counter) Reset() {
	// Prometheus GaugeVec does not have a direct Reset method, so we need to delete all label values
	c.counter.Reset()
}
