package usecase

import (
	"context"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type ReceiverMetricUseCase struct {
	repo ReceiverMetricRepo
}

func NewReceiver(repo ReceiverMetricRepo) *ReceiverMetricUseCase {
	rmu := &ReceiverMetricUseCase{repo: repo}
	return rmu
}

func (receiver *ReceiverMetricUseCase) GetMetricByName(ctx context.Context, name, metricType string) (*metrics2.Metrics, error) {
	var err error
	result := &metrics2.Metrics{ID: name, MType: metricType}
	switch metrics2.MetricType(result.MType) {
	case metrics2.GaugeName:
		gauge, err := receiver.GetGauge(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		result.Value = new(float64)
		result.SetValue(gauge.GetValue())
	case metrics2.CounterName:
		counter, err := receiver.GetCounter(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		result.Delta = new(int64)
		result.SetDelta(counter.GetValue())
	default:
		return nil, metrics2.ErrorUnknownMetricType
	}
	return result, err
}

func (receiver *ReceiverMetricUseCase) GetCounter(ctx context.Context, name string) (*metrics2.Counter, error) {
	return receiver.repo.GetCounter(ctx, name)
}
func (receiver *ReceiverMetricUseCase) GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error) {
	return receiver.repo.GetGauge(ctx, name)
}

func (receiver *ReceiverMetricUseCase) GetMetrics(ctx context.Context) (map[string]string, error) {
	return receiver.repo.GetMetrics(ctx)
}
