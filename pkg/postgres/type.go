package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

const _maxTryConnect = 3

type Postgres struct {
	*sqlx.DB
	isInit bool
}

func New(ctx context.Context, dataSource string) (p *Postgres, err error) {
	p = &Postgres{}
	if dataSource == "" {
		return p, nil
	}
	var conn *sqlx.DB
	for i := 1; i <= _maxTryConnect; i++ {
		conn, err = sqlx.ConnectContext(ctx, "postgres", dataSource)

		if err, ok := err.(*pq.Error); ok && pgerrcode.IsConnectionException(err.Code.Name()) && i < _maxTryConnect {
			time.Sleep(time.Second * time.Duration(i+i-1))
			continue
		}

		if err != nil {
			return p, fmt.Errorf("%w sqlx connect", err)
		}
		break
	}

	p.DB = conn
	p.isInit = true
	return
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
