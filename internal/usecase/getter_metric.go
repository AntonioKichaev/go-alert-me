package usecase

import (
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
)

type ReceiverMetricUseCase struct {
	repo ReceiverMetricRepo
}

func NewReceiver(repo ReceiverMetricRepo) *ReceiverMetricUseCase {
	rmu := &ReceiverMetricUseCase{repo: repo}
	return rmu
}

func (receiver *ReceiverMetricUseCase) GetMetricByName(name, metricType string) (*metrics2.Metrics, error) {
	var err error
	result := &metrics2.Metrics{ID: name, MType: metricType}
	switch metrics2.MetricType(result.MType) {
	case metrics2.GaugeName:
		gauge, err := receiver.GetGauge(result.ID)
		if err != nil {
			return nil, err
		}
		result.Value = new(float64)
		result.SetValue(gauge.GetValue())
	case metrics2.CounterName:
		counter, err := receiver.GetCounter(result.ID)
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

func (receiver *ReceiverMetricUseCase) GetCounter(name string) (*metrics2.Counter, error) {
	return receiver.repo.GetCounter(name)
}
func (receiver *ReceiverMetricUseCase) GetGauge(name string) (*metrics2.Gauge, error) {
	return receiver.repo.GetGauge(name)
}

func (receiver *ReceiverMetricUseCase) GetMetrics() (map[string]string, error) {
	return receiver.repo.GetMetrics()
}
