package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	memstorage "github.com/antoniokichaev/go-alert-me/internal/usecase/repo"
	"github.com/jmoiron/sqlx"
	"strconv"
)

const (
	_countTable = "counts"
	_gaugeTable = "gauges"
)

var _ memstorage.Keeper = New(nil)

type Storage struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func (s *Storage) AddCounter(ctx context.Context, counter *metrics2.Counter) (c *metrics2.Counter, err error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	defer func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = fmt.Errorf("err rollback %w", err)
				return
			}
			return
		}
		err = tx.Commit()
	}()

	sqlReq, args, err := s.builder.
		Select("value").
		From(_countTable).
		Where(sq.Eq{"name": counter.GetName()}).
		ToSql()
	var oldValue int
	err = tx.Get(&oldValue, sqlReq, args...)
	if err != sql.ErrNoRows && err != nil {
		return nil, err
	}

	counter.SetValue(counter.GetValue() + int64(oldValue))
	c, err = s.addCounter(ctx, tx, counter)

	return c, err
}

func (s *Storage) SetGauge(ctx context.Context, gauge *metrics2.Gauge) (g *metrics2.Gauge, err error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	defer func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = fmt.Errorf("err rollback %w", err)
				return
			}
			return
		}
		err = tx.Commit()
	}()
	g, err = s.setGauge(ctx, tx, gauge)
	return g, err
}

func (s *Storage) GetCounter(ctx context.Context, name string) (*metrics2.Counter, error) {
	const fName = "postgres.GetCounter"
	sqlReq, args, err := s.builder.Select("value").From(_countTable).Where(sq.Eq{"name": name}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	m, _ := metrics2.NewCounter(name, 0)
	err = s.db.GetContext(ctx, &m.Value, sqlReq, args...)
	if err != nil {

		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}
	return m, nil
}

func (s *Storage) GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error) {
	const fName = "postgres.GetGauge"
	sqlReq, args, err := s.builder.Select("value").From(_gaugeTable).Where(sq.Eq{"name": name}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	m, _ := metrics2.NewGauge(name, 0.0)
	err = s.db.GetContext(ctx, &m.Value, sqlReq, args...)
	if err != nil {

		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}
	return m, nil
}

func (s *Storage) GetMetrics(ctx context.Context) (map[string]string, error) {
	const fName = "postgres.GetMetrics"

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

func (s *Storage) getCounters(ctx context.Context) (map[string]string, error) {
	const fName = "postgres.getCounters"

	sqlReq, args, err :=
		s.builder.Select("name", "value").From(_countTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	counters := make([]metrics2.Counter, 0)
	err = s.db.SelectContext(ctx, &counters, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Select %w", fName, err)
	}
	mp := make(map[string]string, len(counters))
	for _, c := range counters {
		mp[c.GetName()] = strconv.Itoa(int(c.GetValue()))
	}
	return mp, nil
}

func (s *Storage) getGauges(ctx context.Context) (map[string]string, error) {
	const fName = "postgres.getGauges"

	sqlReq, args, err :=
		s.builder.Select("name", "value").From(_gaugeTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	gauges := make([]metrics2.Gauge, 0)
	err = s.db.SelectContext(ctx, &gauges, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Select %w", fName, err)
	}
	mp := make(map[string]string, len(gauges))
	for _, c := range gauges {
		mp[c.GetName()] = strconv.FormatFloat(c.GetValue(), 'f', -1, 64)
	}
	return mp, nil
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db, builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (s *Storage) addCounter(ctx context.Context, tx *sqlx.Tx, counter *metrics2.Counter) (*metrics2.Counter, error) {
	const fName = "postgres.addCounter"

	sqlReq, args, err := s.builder.
		Insert(_countTable).
		Columns("name", "value").
		Values(counter.GetName(), counter.GetValue()).
		Suffix("ON CONFLICT(name) DO UPDATE SET value=$2").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}

	result, err := tx.ExecContext(ctx, sqlReq, args...)
	if err != nil {

		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s affected %w", fName, err)
	}

	if affected == 0 {
		return nil, fmt.Errorf("%s RowsAffected %w", fName, err)
	}
	return counter, err
}
func (s *Storage) setGauge(ctx context.Context, tx *sqlx.Tx, gauge *metrics2.Gauge) (*metrics2.Gauge, error) {
	const fName = "postgres.setGauge"

	sqlReq, args, err := s.builder.
		Insert(_gaugeTable).
		Columns("name", "value").
		Values(gauge.GetName(), gauge.GetValue()).
		Suffix("ON CONFLICT(name) DO UPDATE SET value =$2").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}

	result, err := tx.ExecContext(ctx, sqlReq, args...)
	if err != nil {

		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s affected %w", fName, err)
	}

	if affected == 0 {
		return nil, fmt.Errorf("%s RowsAffected %w", fName, err)
	}
	return gauge, err
}
