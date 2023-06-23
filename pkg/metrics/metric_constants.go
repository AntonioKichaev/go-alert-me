package metrics

type MetricType string

func (m MetricType) String() string {
	return string(m)
}

const (
	CounterName MetricType = "counter"
	GaugeName   MetricType = "gauge"
)
