package memstorage

import (
	"errors"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
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

func (m *MemStorage) GetCounter(name string) (*metrics2.Counter, error) {
	if val, ok := m.storeCounter[name]; ok {
		return metrics2.NewCounter(name, val)
	}
	return nil, ErrorNotExistMetric
}

func (m *MemStorage) GetGauge(name string) (*metrics2.Gauge, error) {
	if val, ok := m.storeGauge[name]; ok {
		return metrics2.NewGauge(name, val)
	}
	return nil, ErrorNotExistMetric
}

func (m *MemStorage) AddCounter(counter *metrics2.Counter) (*metrics2.Counter, error) {
	m.storeCounter[counter.GetName()] += counter.GetValue()
	counter.SetValue(m.storeCounter[counter.GetName()])
	return counter, nil

}
func (m *MemStorage) SetGauge(gauge *metrics2.Gauge) (*metrics2.Gauge, error) {
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
