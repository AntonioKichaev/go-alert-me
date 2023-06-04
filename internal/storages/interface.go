package storages

//go:generate mockery  --name MetricRepository
type MetricRepository interface {
	GetCounter(metricName string) (int64, error)
	AddCounter(metricName string, value int64)
	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, value float64)
	GetMetrics() map[string]string
}
