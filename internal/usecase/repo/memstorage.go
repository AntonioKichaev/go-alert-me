package memstorage

import (
	"errors"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
	"strconv"
)

var ErrorNotExistMetric = errors.New("doesn't exist metric")

//go:generate mockery --name Keeper
type Keeper interface {
	usecase.UpdaterRepo
	usecase.ReceiverMetricRepo
}

type MemStorage struct {
	storeCounter map[string]int64
	storeGauge   map[string]float64
}

func NewMemStorage() Keeper {
	return newMemStorage()
}
func newMemStorage() *MemStorage {
	return &MemStorage{
		storeCounter: make(map[string]int64, 5),
		storeGauge:   make(map[string]float64, 5),
	}
}

func (m *MemStorage) GetCounter(name string) (*metrics.Counter, error) {
	if val, ok := m.storeCounter[name]; ok {
		return metrics.NewCounter(name, val)
	}
	return nil, ErrorNotExistMetric
}

func (m *MemStorage) GetGauge(name string) (*metrics.Gauge, error) {
	if val, ok := m.storeGauge[name]; ok {
		return metrics.NewGauge(name, val)
	}
	return nil, ErrorNotExistMetric
}

func (m *MemStorage) AddCounter(counter *metrics.Counter) (*metrics.Counter, error) {
	m.storeCounter[counter.GetName()] += counter.GetValue()
	counter.SetValue(m.storeCounter[counter.GetName()])
	return counter, nil

}
func (m *MemStorage) SetGauge(gauge *metrics.Gauge) (*metrics.Gauge, error) {
	m.storeGauge[gauge.GetName()] = gauge.GetValue()

	return gauge, nil
}

func (m *MemStorage) GetMetrics() (map[string]string, error) {
	result := make(map[string]string, len(m.storeGauge)+len(m.storeCounter))
	for key, val := range m.storeGauge {
		result[key] = strconv.FormatFloat(val, 'g', -1, 64)

	}
	for key, val := range m.storeCounter {
		result[key] = strconv.FormatInt(val, 10)
	}
	return result, nil
}
