package postgres

import (
	"context"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/postgres/count"
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/postgres/gauge"
	"github.com/jmoiron/sqlx"
)

var _ memstorage.Keeper = New(nil)

type MetricsCounters interface {
	GetCounters(ctx context.Context) (map[string]string, error)
}
type MetricsGauges interface {
	GetGauges(ctx context.Context) (map[string]string, error)
}

type Storage struct {
	*count.CounterRepo
	*gauge.GaugeRepo
	couterRepo MetricsCounters
	gaugeRepo  MetricsGauges
	db         *sqlx.DB
}

func (s *Storage) GetMetrics(ctx context.Context) (map[string]string, error) {

	mp := make(map[string]string, 0)
	counters, err := s.couterRepo.GetCounters(ctx)
	if err != nil {
		return nil, err
	}
	for key, val := range counters {
		mp[key] = val
	}
	gauges, err := s.gaugeRepo.GetGauges(ctx)
	if err != nil {
		return nil, err
	}
	for key, val := range gauges {
		mp[key] = val
	}
	return mp, nil
}
func (s *Storage) Ping() error {
	return s.db.Ping()
}

func New(db *sqlx.DB) *Storage {
	gauges := gauge.New(db)
	counters := count.New(db)

	st := &Storage{
		db:          db,
		CounterRepo: counters,
		GaugeRepo:   gauges,
		gaugeRepo:   gauges,
		couterRepo:  counters,
	}

	return st
}
