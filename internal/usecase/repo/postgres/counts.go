package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/jmoiron/sqlx"
	"strconv"
)

const _countTable = "counts"

type counterRepo struct {
	builder sq.StatementBuilderType
	db      *sqlx.DB
}

func (cnt *counterRepo) UpdateMetricCounterBatch(ctx context.Context, metrics []metrics2.Counter) error {
	const fName = "postgres.UpdateMetricCounterBatch"

	insertBuilder := cnt.builder.
		Insert(_countTable).
		Columns("name", "value")
	for _, counter := range metrics {
		insertBuilder = insertBuilder.Values(counter.GetName(), counter.GetValue())
	}

	sqlRes, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s builder %w", fName, err)
	}

	result, err := cnt.db.ExecContext(ctx, sqlRes, args...)
	if err != nil {
		return fmt.Errorf("%s ExecContext %w", fName, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s RowsAffected %w", fName, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s RowsAffected 0 %w", fName, err)
	}

	return nil
}
func (cnt *counterRepo) AddCounter(ctx context.Context, counter *metrics2.Counter) (c *metrics2.Counter, err error) {
	c, err = cnt.addCounter(ctx, counter)
	return c, err
}
func (cnt *counterRepo) GetCounter(ctx context.Context, name string) (*metrics2.Counter, error) {
	const fName = "postgres.GetCounter"
	sqlReq, args, err := cnt.builder.Select("sum(value)").From(_countTable).GroupBy("name").Where(sq.Eq{"name": name}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	m := &metrics2.Counter{Name: name}
	err = cnt.db.GetContext(ctx, &m.Value, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}
	return m, nil
}
func (cnt *counterRepo) getCounters(ctx context.Context) (map[string]string, error) {
	const fName = "postgres.getCounters"

	sqlReq, args, err :=
		cnt.builder.Select("name", "sum(value) as value").
			GroupBy("name").
			From(_countTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	counters := make([]metrics2.Counter, 0)
	err = cnt.db.SelectContext(ctx, &counters, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Select %w", fName, err)
	}
	mp := make(map[string]string, len(counters))
	for _, c := range counters {
		mp[c.GetName()] = strconv.Itoa(int(c.GetValue()))
	}
	return mp, nil
}
func (cnt *counterRepo) addCounter(ctx context.Context, counter *metrics2.Counter) (*metrics2.Counter, error) {
	const fName = "postgres.addCounter"

	sqlReq, args, err := cnt.builder.
		Insert(_countTable).
		Columns("name", "value").
		Values(counter.GetName(), counter.GetValue()).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}

	result, err := cnt.db.ExecContext(ctx, sqlReq, args...)
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
