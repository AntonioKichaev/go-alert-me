package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/jmoiron/sqlx"
)

var _ memstorage.Keeper = New(nil)

type Storage struct {
	gaugeRepo
	counterRepo
	db *sqlx.DB
}

func (s *Storage) GetMetrics(ctx context.Context) (map[string]string, error) {

	mp := make(map[string]string, 0)
	counters, err := s.getCounters(ctx)
	if err != nil {
		return nil, err
	}
	for key, val := range counters {
		mp[key] = val
	}
	gauges, err := s.getGauges(ctx)
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
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	st := &Storage{
		db:          db,
		gaugeRepo:   gaugeRepo{db: db, builder: builder},
		counterRepo: counterRepo{db: db, builder: builder},
	}

	return st
}
