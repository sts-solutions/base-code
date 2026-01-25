package ccmetrics

type DatabaseMetricsHandler interface {
	RequestDuration(source string, duration float64)
	RequestCounterInc(source string)
	ErrorInc(source string)
	UpdateMaxOpenConnections(count int)
	UpdateOpenConnections(count int)
	UpdateWaitCount(count int)
	UpdateWaitDuration(totalDuration int)
	UpdateIdleCount(count int)
	UpdateInUseCount(count int)
	QueueCount(name string, count int)
}
