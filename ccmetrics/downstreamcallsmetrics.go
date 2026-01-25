package ccmetrics

import "google.golang.org/grpc"

type DownstreamCallsMetricsHandler interface {
	HttpRequestDuration(recipient, method, path, code string, duration float64)
	HttpRequestCounterInc(recipient, method, path, code string)
	HttpErrorInc(recipient, method, path, code string)
}

type DownstreamGrpcCallsMetricsHandler interface {
	GrpcRequestDuration(recipient, path, code string, duration float64)
	GrpcRequestCounterInc(recipient, path, code string)
	GrpcErrorInc(recipient, method, code string)
	Interceptor() grpc.UnaryClientInterceptor
}
