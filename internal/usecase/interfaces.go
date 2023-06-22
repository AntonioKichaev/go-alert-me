package usecase

import "github.com/antoniokichaev/go-alert-me/internal/entity"

type (
	Updater interface {
		AddCounter(name string, value any) error
		SetGauge(name string, value any) error
	}
	UpdaterRepo interface {
		AddCounter(counter *entity.Counter) error
		SetGauge(gauge *entity.Gauge) error
	}
)

type (
	ReceiverMetric interface {
		GetCounter(name string) (*entity.Counter, error)
		GetGauge(name string) (*entity.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
	ReceiverMetricRepo interface {
		GetCounter(name string) (*entity.Counter, error)
		GetGauge(name string) (*entity.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
)
