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
	return mm.mem.GetMetrics()
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
	return mm.mem.GetCounter(metricName) // todo:error if metricName doesn't exist
}
func (mm *MockMemStorage) GetGauge(metricName string) (float64, error) {
	return mm.mem.GetGauge(metricName)
}

func TestMemStorage_AddCounter(t *testing.T) {
	tt := map[string]struct {
		metricName  string
		values      []int64
		wantMetrics map[string]string
		want        int64
		wantErr     error
	}{
		"one_value": {
			metricName:  "ram",
			values:      []int64{352},
			want:        352,
			wantMetrics: map[string]string{"ram": "352"},
			wantErr:     nil,
		},
		"many_values": {
			metricName:  "koef",
			values:      []int64{1, 2, 3, 4, 5},
			want:        15,
			wantMetrics: map[string]string{"koef": "15"},
			wantErr:     nil,
		},
		"many_values_negative_include": {
			metricName:  "ram",
			values:      []int64{1, 2, 3, 4, -5},
			want:        5,
			wantMetrics: map[string]string{"ram": "5"},
			wantErr:     nil,
		},
		"error not found": {
			metricName:  "ram",
			values:      []int64{},
			want:        0,
			wantMetrics: map[string]string{},
			wantErr:     ErrorNotExistMetric,
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
			req.ErrorIs(tc.wantErr, err)
			req.EqualValues(tc.wantMetrics, storage.GetMetrics(), "GetMetrics()")

		})
	}

}

func TestMemStorage_SetGauge(t *testing.T) {
	tt := map[string]struct {
		metricName  string
		values      []float64
		want        float64
		wantMetrics map[string]string
		wantErr     error
	}{
		"one_value": {
			metricName:  "ram",
			values:      []float64{352, 0, 2},
			want:        2,
			wantMetrics: map[string]string{"ram": "2"},
			wantErr:     nil,
		},
		"many_values": {
			metricName:  "ram",
			values:      []float64{1, 2, 3, 4, 5},
			want:        5,
			wantMetrics: map[string]string{"ram": "5"},
			wantErr:     nil,
		},
		"many_values_negative_include": {
			metricName:  "ram",
			values:      []float64{1, 2, 3, 4, -5},
			want:        -5,
			wantMetrics: map[string]string{"ram": "-5"},
			wantErr:     nil,
		},
		"ErrorNotExist": {
			metricName:  "ram",
			values:      []float64{},
			want:        0,
			wantMetrics: map[string]string{},
			wantErr:     ErrorNotExistMetric,
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
			req.ErrorIs(tc.wantErr, err)
			req.EqualValues(tc.wantMetrics, storage.GetMetrics(), "GetMetrics()")

		})
	}

}
