package gauge

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	metrics2 "github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/jmoiron/sqlx"
	"strconv"
)

const _gaugeTable = "gauges"

type GaugeRepo struct {
	builder sq.StatementBuilderType
	db      *sqlx.DB
}

func (grp *GaugeRepo) UpdateMetricGaugeBatch(ctx context.Context, metrics []metrics2.Gauge) error {
	const fName = "postgres.UpdateMetricCounterBatch"
	tx, err := grp.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s BeginTxx %w", fName, err)
	}

	defer func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = fmt.Errorf("err rollback %w", err)
				return
			}
		}
	}()

	insertBuilder := grp.builder.
		Insert(_gaugeTable).
		Columns("name", "value")
	for _, counter := range metrics {
		insertBuilder = insertBuilder.Values(counter.GetName(), counter.GetValue())
	}

	sqlRes, args, err := insertBuilder.Suffix("ON CONFLICT(name) DO UPDATE SET value=EXCLUDED.value").ToSql()
	if err != nil {
		return fmt.Errorf("%s builder %w", fName, err)
	}

	result, err := tx.ExecContext(ctx, sqlRes, args...)
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
	return tx.Commit()
}
func (grp *GaugeRepo) SetGauge(ctx context.Context, gauge *metrics2.Gauge) (g *metrics2.Gauge, err error) {
	tx, err := grp.db.BeginTxx(ctx, nil)
	defer func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = fmt.Errorf("err rollback %w", err)
				return
			}
		}

	}()
	g, err = grp.setGauge(ctx, tx, gauge)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return g, err
}
func (grp *GaugeRepo) GetGauge(ctx context.Context, name string) (*metrics2.Gauge, error) {
	const fName = "postgres.GetGauge"
	sqlReq, args, err := grp.builder.Select("value").From(_gaugeTable).Where(sq.Eq{"name": name}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	m, _ := metrics2.NewGauge(name, 0.0)
	err = grp.db.GetContext(ctx, &m.Value, sqlReq, args...)
	if err != nil {

		return nil, fmt.Errorf("%s Exec %w", fName, err)
	}
	return m, nil
}
func (grp *GaugeRepo) GetGauges(ctx context.Context) (map[string]string, error) {
	const fName = "postgres.GetGauges"

	sqlReq, args, err :=
		grp.builder.Select("name", "value").From(_gaugeTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s builder %w", fName, err)
	}
	gauges := make([]metrics2.Gauge, 0)
	err = grp.db.SelectContext(ctx, &gauges, sqlReq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s Select %w", fName, err)
	}
	mp := make(map[string]string, len(gauges))
	for _, c := range gauges {
		mp[c.GetName()] = strconv.FormatFloat(c.GetValue(), 'f', -1, 64)
	}
	return mp, nil
}
func (grp *GaugeRepo) setGauge(ctx context.Context, tx *sqlx.Tx, gauge *metrics2.Gauge) (*metrics2.Gauge, error) {
	const fName = "postgres.setGauge"

	sqlReq, args, err := grp.builder.
		Insert(_gaugeTable).
		Columns("name", "value").
		Values(gauge.GetName(), gauge.GetValue()).
		Suffix("ON CONFLICT(name) DO UPDATE SET value =EXCLUDED.value").
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

func New(db *sqlx.DB) *GaugeRepo {
	return &GaugeRepo{
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		db:      db,
	}
}
