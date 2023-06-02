package storages

type Keeper interface {
	AddCounter(metricName string, value int64)
	SetGauge(metricName string, value float64)
}
