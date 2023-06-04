package memstorage

import (
	"github.com/antoniokichaev/go-alert-me/internal/storages"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockKeeper interface {
	storages.MetricRepository
}

type MockMemStorage struct {
	mem *MemStorage
}

func (mm *MockMemStorage) GetMetrics() map[string]string {
	return nil

}

func NewMockMemStorage() MockKeeper {
	return &MockMemStorage{mem: newMemStorage()}
}
func (mm *MockMemStorage) AddCounter(metricName string, value int64) {
	mm.mem.AddCounter(metricName, value)
}
func (mm *MockMemStorage) SetGauge(metricName string, value float64) {
	mm.mem.SetGauge(metricName, value)
}

func (mm *MockMemStorage) GetCounter(metricName string) (int64, error) {
	return mm.mem.storeCounter[metricName], nil // todo:error if metricName doesn't exist
}
func (mm *MockMemStorage) GetGauge(metricName string) (float64, error) {
	return mm.mem.storeGauge[metricName], nil
}

func TestMemStorage_AddCounter(t *testing.T) {
	tt := map[string]struct {
		metricName string
		values     []int64
		want       int64
	}{
		"one_value": {
			metricName: "ram",
			values:     []int64{352},
			want:       352,
		},
		"many_values": {
			metricName: "ram",
			values:     []int64{1, 2, 3, 4, 5},
			want:       15,
		},
		"many_values_negative_include": {
			metricName: "ram",
			values:     []int64{1, 2, 3, 4, -5},
			want:       5,
		},
	}

	req := require.New(t)
	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			storage := NewMockMemStorage()
			for _, val := range tc.values {
				storage.AddCounter(tc.metricName, val)
			}
			got, err := storage.GetCounter(tc.metricName)
			req.EqualValues(tc.want, got)
			req.NoError(err)

		})
	}

}

func TestMemStorage_SetGauge(t *testing.T) {
	tt := map[string]struct {
		metricName string
		values     []float64
		want       float64
	}{
		"one_value": {
			metricName: "ram",
			values:     []float64{352, 0, 2},
			want:       2,
		},
		"many_values": {
			metricName: "ram",
			values:     []float64{1, 2, 3, 4, 5},
			want:       5,
		},
		"many_values_negative_include": {
			metricName: "ram",
			values:     []float64{1, 2, 3, 4, -5},
			want:       -5,
		},
	}
	req := require.New(t)
	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			storage := NewMockMemStorage()
			for _, val := range tc.values {
				storage.SetGauge(tc.metricName, val)
			}
			got, err := storage.GetGauge(tc.metricName)
			req.EqualValues(tc.want, got)
			req.NoError(err)

		})
	}

}
