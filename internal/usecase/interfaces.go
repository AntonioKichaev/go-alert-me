package usecase

import (
	"context"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type (
	Updater interface {
		UpdateMetricByParams(ctx context.Context, name, metricType string, value any) (*metrics2.Metrics, error)
		UpdateMetric(ctx context.Context, m *metrics2.Metrics) (*metrics2.Metrics, error)
		AddCounter(ctx context.Context, name string, value any) (*metrics2.Counter, error)
		SetGauge(ctx context.Context, name string, value any) (*metrics2.Gauge, error)
		UpdateMetricBatch(ctx context.Context, metrics []metrics2.Metrics) error
	}
	UpdaterRepo interface {
		AddCounter(ctx context.Context, counter *metrics2.Counter) (*metrics2.Counter, error)
		SetGauge(ctx context.Context, gauge *metrics2.Gauge) (*metrics2.Gauge, error)
		UpdateMetricCounterBatch(ctx context.Context, metrics []metrics2.Counter) error
		UpdateMetricGaugeBatch(ctx context.Context, metrics []metrics2.Gauge) error
	}
)

type (
	ReceiverMetric interface {
		GetMetricByName(ctx context.Context, name, metricType string) (*metrics2.Metrics, error)
		GetCounter(ctx context.Context, name string) (*metrics2.Counter, error)
		GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error)
		GetMetrics(ctx context.Context) (map[string]string, error)
	}
	ReceiverMetricRepo interface {
		GetCounter(ctx context.Context, name string) (*metrics2.Counter, error)
		GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error)
		GetMetrics(ctx context.Context) (map[string]string, error)
	}
)
