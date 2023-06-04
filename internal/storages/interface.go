package storages

//go:generate mockery  --name MetricRepository
type MetricRepository interface {
	AddCounter(metricName string, value int64)
	SetGauge(metricName string, value float64)
}
