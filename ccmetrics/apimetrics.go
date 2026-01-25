package ccmetrics

type HttpApiMetricsHandler interface {
	RequestDuration(method string, path string, code string, duration float64)
	RequestCounterInc(method string, path string, code string)
	ErrorInc(method string, path string, code string)
}

// GrpcApiMetricsHandler
type GrpcApiMetricsHandler interface {
	RequestDuration(path string, code string, duration float64)
	RequestCounterInc(path string, code string)
	ErrorInc(path string, code string)
}
