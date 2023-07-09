package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	*sqlx.DB
	isInit bool
}

func New(ctx context.Context, dataSource string) (*Postgres, error) {
	p := &Postgres{}
	if dataSource == "" {
		return p, nil
	}

	conn, err := sqlx.ConnectContext(ctx, "postgres", dataSource)
	if err != nil {
		return p, fmt.Errorf("%w sqlx connext", err)
	}
	p.DB = conn
	p.isInit = true
	return p, err
}

func (p *Postgres) Ping() error {
	if !p.isInit {
		return errors.New("doesn't init")
	}
	return p.DB.Ping()
}
func (p *Postgres) Close() error {
	return p.DB.Close()
}
