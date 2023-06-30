package usecase

import (
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type UpdaterUseCase struct {
	repo UpdaterRepo
}

func NewUpdater(repo UpdaterRepo) *UpdaterUseCase {
	return &UpdaterUseCase{repo: repo}
}

func (u *UpdaterUseCase) AddCounter(name string, value any) (*metrics2.Counter, error) {
	c, err := metrics2.NewCounter(name, value)
	if err != nil {
		return nil, err
	}
	return u.repo.AddCounter(c)
}

func (u *UpdaterUseCase) SetGauge(name string, value any) (*metrics2.Gauge, error) {
	g, err := metrics2.NewGauge(name, value)
	if err != nil {
		return nil, err
	}
	g, err = u.repo.SetGauge(g)
	return g, err
}
func (u *UpdaterUseCase) UpdateMetric(m *metrics2.Metrics) (*metrics2.Metrics, error) {
	return u.updateMetric(m)
}

func (u *UpdaterUseCase) UpdateMetricByParams(name, metricType string, value any) (*metrics2.Metrics, error) {
	m, err := metrics2.NewMetrics(
		metrics2.SetName(name),
		metrics2.SetMetricType(metricType),
		metrics2.SetValueOrDelta(value))
	if err != nil {
		return nil, err
	}

	return u.updateMetric(m)
}
func (u *UpdaterUseCase) updateMetric(m *metrics2.Metrics) (*metrics2.Metrics, error) {
	var err error
	switch metrics2.MetricType(m.MType) {
	case metrics2.GaugeName:
		g, err := m.ToGauge()
		if err != nil {
			return nil, err
		}
		g, err = u.repo.SetGauge(g)
		if err != nil {
			m.SetValue(g.GetValue())
		}
	case metrics2.CounterName:
		c, err := m.ToCounter()
		if err != nil {
			return nil, err
		}
		c, err = u.repo.AddCounter(c)
		if err != nil {
			return nil, err
		}
		m.SetDelta(c.GetValue())
	default:
		err = metrics2.ErrorUnknownMetricType
	}
	return m, err
}
