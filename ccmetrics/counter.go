package ccmetrics

type Counter interface {
	Inc(labels ...string)
	Set(value float64, labels ...string)
	Reset()
}
