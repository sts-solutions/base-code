package ccmetrics

type ApplicationMetricsHandler interface {
	VersionInfo(version, goVersion string)
	ErrorInc(processName string)
	PanicInc(correlationID string)
}
