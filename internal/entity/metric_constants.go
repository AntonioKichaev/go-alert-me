package entity

type MetricType string

const (
	CounterName MetricType = "counter"
	GaugeName   MetricType = "gauge"
)
