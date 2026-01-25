package ccmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sts-solutions/base-code/ccmetrics"
)

type ApplicationMetricsHandler struct {
	verion       *prometheus.GaugeVec
	errorCounter *prometheus.CounterVec
	panicCounter *prometheus.CounterVec
}

func NewApplicationMetricsHandler(namespace string) ccmetrics.ApplicationMetricsHandler {
	amh := &ApplicationMetricsHandler{}

	amh.verion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "version_info",
			Help:      "Application version info",
		},
		[]string{"version", "go_version"},
	)

	amh.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      " errors_total",
			Help:      "Total number of errors",
		},
		[]string{"process_name"},
	)

	amh.panicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "panics_total",
			Help:      "Total number of panics",
		},
		[]string{"correlation_id"},
	)

	return amh
}

func (amh *ApplicationMetricsHandler) VersionInfo(version, goVersion string) {
	amh.verion.WithLabelValues(version, goVersion).Set(1)
}

func (amh *ApplicationMetricsHandler) ErrorInc(processName string) {
	amh.errorCounter.WithLabelValues(processName).Inc()
}

func (amh *ApplicationMetricsHandler) PanicInc(correlationID string) {
	amh.panicCounter.WithLabelValues(correlationID).Inc()
}
