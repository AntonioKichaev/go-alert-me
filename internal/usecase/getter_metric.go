package usecase

import "github.com/antoniokichaev/go-alert-me/internal/entity"

type ReceiverMetricUseCase struct {
	repo ReceiverMetricRepo
}

func NewReceiver(repo ReceiverMetricRepo) *ReceiverMetricUseCase {
	rmu := &ReceiverMetricUseCase{repo: repo}
	return rmu
}

func (receiver *ReceiverMetricUseCase) GetCounter(name string) (*entity.Counter, error) {
	return receiver.repo.GetCounter(name)
}
func (receiver *ReceiverMetricUseCase) GetGauge(name string) (*entity.Gauge, error) {
	return receiver.repo.GetGauge(name)
}

func (receiver *ReceiverMetricUseCase) GetMetrics() (map[string]string, error) {
	return receiver.repo.GetMetrics()
}
