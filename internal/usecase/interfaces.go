package usecase

import (
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
)

type (
	Updater interface {
		AddCounter(name string, value any) (*metrics.Counter, error)
		SetGauge(name string, value any) (*metrics.Gauge, error)
	}
	UpdaterRepo interface {
		AddCounter(counter *metrics.Counter) (*metrics.Counter, error)
		SetGauge(gauge *metrics.Gauge) (*metrics.Gauge, error)
	}
)

type (
	ReceiverMetric interface {
		GetCounter(name string) (*metrics.Counter, error)
		GetGauge(name string) (*metrics.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
	ReceiverMetricRepo interface {
		GetCounter(name string) (*metrics.Counter, error)
		GetGauge(name string) (*metrics.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
)
