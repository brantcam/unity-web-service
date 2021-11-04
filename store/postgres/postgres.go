package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/qustavo/dotsql"
	"github.com/jmoiron/sqlx"

	// postgres driver for the database/sql package
	_ "github.com/lib/pq"
)

type Conn struct {
	Queries *dotsql.DotSql
	Client *sqlx.DB

	retry        int
	retryBackoff time.Duration
}

func New(ctx context.Context, c Config) (*Conn, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host,
			c.Port,
			c.Username,
			c.Password,
			c.Name,
			"disable",
		),
	)
	if err != nil {
		return nil, err
	}
	dot, err := dotsql.LoadFromFile("./store/postgres/queries.sql")
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		Queries: dot,
		Client: db,
		retry: 3,
		retryBackoff: 5 * time.Second,
	}

	return conn, db.PingContext(ctx)
}

func (c *Conn) Health(ctx context.Context) error {
	var err error

	for i := 0; i <= c.retry; i++ {
		if err = c.Client.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(c.retryBackoff)
	}

	return err
}
