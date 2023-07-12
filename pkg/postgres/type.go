package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const _maxTryConnect = 3

type Postgres struct {
	*sqlx.DB
	isInit bool
}

func New(ctx context.Context, dataSource string) (pg *Postgres, err error) {
	pg = &Postgres{}
	if dataSource == "" {
		return pg, errors.New("dataSource is empty")
	}
	var conn *sqlx.DB
	for i := 1; i <= _maxTryConnect; i++ {
		conn, err = sqlx.ConnectContext(ctx, "postgres", dataSource)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+i-1))
	}

	if err != nil {
		return pg, fmt.Errorf("%w sqlx connect", err)
	}

	pg.DB = conn
	pg.isInit = true
	return
}

func (p *Postgres) Ping() error {
	if !p.isInit {
		return errors.New("doesn't init")
	}
	return p.DB.Ping()
}
func (p *Postgres) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}
