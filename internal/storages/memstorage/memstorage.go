package memstorage

import (
	"github.com/antoniokichaev/go-alert-me/internal/storages"
)

type MemStorage struct {
	storeCounter map[string]int64
	storeGauge   map[string]float64
}

func NewMemStorage() storages.Keeper {
	return newMemStorage()
}
func newMemStorage() *MemStorage {
	return &MemStorage{
		storeCounter: make(map[string]int64, 5),
		storeGauge:   make(map[string]float64, 5),
	}
}

func (m *MemStorage) AddCounter(metricName string, value int64) {
	m.storeCounter[metricName] += value

}
func (m *MemStorage) SetGauge(metricName string, value float64) {
	m.storeGauge[metricName] = value
}
