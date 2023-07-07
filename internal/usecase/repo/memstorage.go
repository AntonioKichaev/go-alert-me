package memstorage

import (
	"context"
	"errors"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/pkg/memorystorage"
	"strconv"
)

//go:generate mockery --name Keeper
type Keeper interface {
	usecase.UpdaterRepo
	usecase.ReceiverMetricRepo
	Ping() error
}

type MemStorage struct {
	storeCounter *memorystorage.MemoryStorage
	storeGauge   *memorystorage.MemoryStorage
}

func NewMemStorage(storeCounter, storeGauge *memorystorage.MemoryStorage) Keeper {
	return newMemStorage(storeCounter, storeGauge)
}
func newMemStorage(storeCounter, storeGauge *memorystorage.MemoryStorage) *MemStorage {
	return &MemStorage{
		storeCounter: storeCounter,
		storeGauge:   storeGauge,
	}
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (*metrics2.Counter, error) {
	if val, err := m.storeCounter.Get(name); err != nil {
		return nil, err
	} else {
		return metrics2.NewCounter(name, val)
	}
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error) {
	if val, err := m.storeGauge.Get(name); err != nil {
		return nil, err
	} else {
		return metrics2.NewGauge(name, val)
	}
}

func (m *MemStorage) AddCounter(ctx context.Context, counter *metrics2.Counter) (*metrics2.Counter, error) {
	old, err := m.GetCounter(ctx, counter.GetName())
	if errors.Is(err, memorystorage.ErrorNotExistMetric) {
		old, _ = metrics2.NewCounter(counter.GetName(), 0)

	} else if err != nil {
		return nil, err
	}
	counter.SetValue(old.GetValue() + counter.GetValue())
	err = m.storeCounter.Set(counter.GetName(), strconv.FormatInt(counter.GetValue(), 10))
	return counter, err

}
func (m *MemStorage) SetGauge(ctx context.Context, gauge *metrics2.Gauge) (*metrics2.Gauge, error) {
	err := m.storeGauge.Set(gauge.GetName(), strconv.FormatFloat(gauge.GetValue(), 'f', -1, 64))
	return gauge, err
}

func (m *MemStorage) GetMetrics(ctx context.Context) (map[string]string, error) {
	g := m.storeGauge.GetDump()
	c := m.storeCounter.GetDump()
	result := make(map[string]string, len(g)+len(c))
	for key, val := range g {
		result[key] = val
	}
	for key, val := range c {
		result[key] = val
	}
	return result, nil
}
func (m *MemStorage) Ping() error {
	return errors.New("doesn't impliment")
}
