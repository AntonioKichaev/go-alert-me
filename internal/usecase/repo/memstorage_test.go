package memstorage

import (
	"github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorage_AddCounter(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	tt := map[string]struct {
		fields      []fields
		wantMetrics map[string]string
		want        int64
		wantErr     error
	}{
		"one_value": {
			fields:      []fields{{"ram", 352}},
			want:        352,
			wantMetrics: map[string]string{"ram": "352"},
			wantErr:     nil,
		},
		"many_values": {
			fields: []fields{
				{"koef", 1},
				{"koef", 2},
				{"koef", 3},
				{"koef", 4},
				{"koef", 5},
			},
			want:        15,
			wantMetrics: map[string]string{"koef": "15"},
			wantErr:     nil,
		},
		"many_values_negative_include": {
			fields: []fields{
				{"ram", 1},
				{"ram", 2},
				{"ram", 3},
				{"ram", 4},
				{"ram", -5},
			},
			want:        5,
			wantMetrics: map[string]string{"ram": "5"},
			wantErr:     nil,
		},
	}

	req := require.New(t)
	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			storage := mocks.NewKeeper(t)
			res := int64(0)
			for _, val := range tc.fields {
				c, err := metrics.NewCounter(val.name, val.value)
				assert.NoError(t, err, "create newCounter err")
				storage.EXPECT().AddCounter(c).Return(nil, tc.wantErr)
				res += c.GetValue()
				_, _ = storage.AddCounter(c)
			}
			metricName := ""
			if len(tc.fields) != 0 {
				metricName = tc.fields[0].name
			}
			storage.EXPECT().GetCounter(metricName).Return(&metrics.Counter{Name: metricName, Value: res}, tc.wantErr)
			got, err := storage.GetCounter(metricName)
			req.EqualValues(tc.want, got.GetValue())
			req.ErrorIs(tc.wantErr, err)
			storage.EXPECT().GetMetrics().Return(tc.wantMetrics, tc.wantErr)
			metrics, err := storage.GetMetrics()
			req.NoError(err)
			req.EqualValues(tc.wantMetrics, metrics, "GetMetrics()")

		})
	}

}

//func TestMemStorage_SetGauge(t *testing.T) {
//	tt := map[string]struct {
//		metricName  string
//		values      []float64
//		want        float64
//		wantMetrics map[string]string
//		wantErr     error
//	}{
//		"one_value": {
//			metricName:  "ram",
//			values:      []float64{352, 0, 2},
//			want:        2,
//			wantMetrics: map[string]string{"ram": "2"},
//			wantErr:     nil,
//		},
//		"many_values": {
//			metricName:  "ram",
//			values:      []float64{1, 2, 3, 4, 5},
//			want:        5,
//			wantMetrics: map[string]string{"ram": "5"},
//			wantErr:     nil,
//		},
//		"many_values_negative_include": {
//			metricName:  "ram",
//			values:      []float64{1, 2, 3, 4, -5},
//			want:        -5,
//			wantMetrics: map[string]string{"ram": "-5"},
//			wantErr:     nil,
//		},
//		"ErrorNotExist": {
//			metricName:  "ram",
//			values:      []float64{},
//			want:        0,
//			wantMetrics: map[string]string{},
//			wantErr:     ErrorNotExistMetric,
//		},
//	}
//	req := require.New(t)
//	for key, tc := range tt {
//		t.Run(key, func(t *testing.T) {
//			storage := mocks.NewKeeper(t)
//			for _, val := range tc.values {
//				storage.SetGauge(tc.metricName, val)
//			}
//			got, err := storage.GetGauge(tc.metricName)
//			req.EqualValues(tc.want, got)
//			req.ErrorIs(tc.wantErr, err)
//			req.EqualValues(tc.wantMetrics, storage.GetMetrics(), "GetMetrics()")
//
//		})
//	}
//
//}
