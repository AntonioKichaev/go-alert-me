package usecase

import (
	"context"
	"fmt"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type UpdaterUseCase struct {
	repo UpdaterRepo
}

func NewUpdater(repo UpdaterRepo) *UpdaterUseCase {
	return &UpdaterUseCase{repo: repo}
}

func (u *UpdaterUseCase) AddCounter(ctx context.Context, name string, value any) (*metrics2.Counter, error) {
	c, err := metrics2.NewCounter(name, value)
	if err != nil {
		return nil, err
	}
	return u.repo.AddCounter(ctx, c)
}

func (u *UpdaterUseCase) SetGauge(ctx context.Context, name string, value any) (*metrics2.Gauge, error) {
	g, err := metrics2.NewGauge(name, value)
	if err != nil {
		return nil, err
	}
	g, err = u.repo.SetGauge(ctx, g)
	return g, err
}
func (u *UpdaterUseCase) UpdateMetric(ctx context.Context, m *metrics2.Metrics) (*metrics2.Metrics, error) {
	return u.updateMetric(ctx, m)
}

func (u *UpdaterUseCase) UpdateMetricByParams(ctx context.Context, name, metricType string, value any) (*metrics2.Metrics, error) {
	m, err := metrics2.NewMetrics(
		metrics2.SetName(name),
		metrics2.SetMetricType(metricType),
		metrics2.SetValueOrDelta(value))
	if err != nil {
		return nil, err
	}

	return u.updateMetric(ctx, m)
}
func (u *UpdaterUseCase) updateMetric(ctx context.Context, m *metrics2.Metrics) (*metrics2.Metrics, error) {
	var err error
	switch metrics2.MetricType(m.MType) {
	case metrics2.GaugeName:
		g, err := m.ToGauge()
		if err != nil {
			return nil, err
		}
		g, err = u.repo.SetGauge(ctx, g)
		if err != nil {
			m.SetValue(g.GetValue())
		}
	case metrics2.CounterName:
		c, err := m.ToCounter()
		if err != nil {
			return nil, err
		}
		c, err = u.repo.AddCounter(ctx, c)
		if err != nil {
			return nil, err
		}
		m.SetDelta(c.GetValue())
	default:
		err = metrics2.ErrorUnknownMetricType
	}
	return m, err
}

func (u *UpdaterUseCase) UpdateMetricBatch(ctx context.Context, metrics []metrics2.Metrics) error {
	const fName = "UpdaterUseCase.UpdateMetricBatch"

	counters := make([]metrics2.Counter, 0, len(metrics))
	gaugesUniq := make(map[string]metrics2.Gauge, len(metrics))

	for _, m := range metrics {
		switch metrics2.MetricType(m.MType) {
		case metrics2.GaugeName:

			g, err := m.ToGauge()
			if err != nil {
				return fmt.Errorf("%s ToGauge %w", fName, err)
			}
			gaugesUniq[g.GetName()] = *g
		case metrics2.CounterName:
			c, err := m.ToCounter()
			if err != nil {
				return fmt.Errorf("%s ToCounter %w", fName, err)
			}
			counters = append(counters, *c)
		default:
			return metrics2.ErrorUnknownMetricType
		}
	}
	gauges := make([]metrics2.Gauge, 0, len(gaugesUniq))
	for _, g := range gaugesUniq {
		gauges = append(gauges, g)
	}
	if len(counters) > 0 {
		err := u.repo.UpdateMetricCounterBatch(ctx, counters)
		if err != nil {
			return fmt.Errorf("%s UpdateMetricCounterBatch %w", fName, err)
		}
	}

	if len(gauges) > 0 {
		err := u.repo.UpdateMetricGaugeBatch(ctx, gauges)
		if err != nil {
			return fmt.Errorf("%s UpdateMetricGaugeBatch %w", fName, err)
		}
	}

	return nil
}
