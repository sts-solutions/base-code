package ccmetrics

const (
	HttpInSubsystem  = "http"
	GrpcInSubsystem  = "grpc"
	HttpOutSubsystem = "http_outgoing"
	GrpcOutSubsystem = "grpc_outgoing"
)

var (
	DefaultBuckets           = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 25, 30}
	GrpcClientDefaultBuckets = []float64{.001, 0.003, .005, .01, .025, .05, .1, .25, .5, 1, 2.5}
)
