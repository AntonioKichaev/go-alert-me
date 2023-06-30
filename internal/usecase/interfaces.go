package usecase

import (
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type (
	Updater interface {
		UpdateMetricByParams(name, metricType string, value any) (*metrics2.Metrics, error)
		UpdateMetric(*metrics2.Metrics) (*metrics2.Metrics, error)
		AddCounter(name string, value any) (*metrics2.Counter, error)
		SetGauge(name string, value any) (*metrics2.Gauge, error)
	}
	UpdaterRepo interface {
		AddCounter(counter *metrics2.Counter) (*metrics2.Counter, error)
		SetGauge(gauge *metrics2.Gauge) (*metrics2.Gauge, error)
	}
)

type (
	ReceiverMetric interface {
		GetMetricByName(name, metricType string) (*metrics2.Metrics, error)
		GetCounter(name string) (*metrics2.Counter, error)
		GetGauge(name string) (*metrics2.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
	ReceiverMetricRepo interface {
		GetCounter(name string) (*metrics2.Counter, error)
		GetGauge(name string) (*metrics2.Gauge, error)
		GetMetrics() (map[string]string, error)
	}
)
