package metrics

type Option func(m *Metrics)

func SetName(name string) Option {
	return func(m *Metrics) {
		m.ID = name
	}
}
func SetMetricType(mType string) Option {
	return func(m *Metrics) {
		m.MType = mType
	}
}

func SetValueOrDelta(value any) Option {
	return func(m *Metrics) {
		m.tmpValue = value

	}
}
