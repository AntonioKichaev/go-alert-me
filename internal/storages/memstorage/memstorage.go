package memstorage

import (
	"errors"
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"strconv"
)

var ErrorNotExistMetric = errors.New("doesn't exist metric")

type MemStorage struct {
	storeCounter map[string]int64
	storeGauge   map[string]float64
}

func NewMemStorage() storages.MetricRepository {
	return newMemStorage()
}
func newMemStorage() *MemStorage {
	return &MemStorage{
		storeCounter: make(map[string]int64, 5),
		storeGauge:   make(map[string]float64, 5),
	}
}

func (m *MemStorage) GetCounter(metricName string) (int64, error) {
	if val, ok := m.storeCounter[metricName]; ok {
		return val, nil
	}
	return 0, ErrorNotExistMetric
}

func (m *MemStorage) GetGauge(metricName string) (float64, error) {
	if val, ok := m.storeGauge[metricName]; ok {
		return val, nil
	}
	return 0, ErrorNotExistMetric
}

func (m *MemStorage) AddCounter(metricName string, value int64) {
	m.storeCounter[metricName] += value

}
func (m *MemStorage) SetGauge(metricName string, value float64) {
	m.storeGauge[metricName] = value
}

func (m *MemStorage) GetMetrics() map[string]string {
	result := make(map[string]string, len(m.storeGauge)+len(m.storeCounter))
	for key, val := range m.storeGauge {
		result[key] = strconv.FormatFloat(val, 'g', -1, 64)

	}
	for key, val := range m.storeCounter {
		result[key] = strconv.FormatInt(val, 10)
	}
	return result
}
