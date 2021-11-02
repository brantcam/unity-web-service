package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	// postgres driver for the database/sql package
	_ "github.com/lib/pq"
)

type Conn struct {
	client  *sqlx.DB

	retry int
	retryBackoff time.Duration
}

func New(ctx context.Context, c Config) (*Conn, error) {
	var conn Conn

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
	conn.client = db
	conn.retry = 3 // todo: put these into a config and read it in
	conn.retryBackoff = 5 * time.Second // todo: put these into a config and read it in

	return &conn, db.PingContext(ctx)
}

func (c *Conn) Health(ctx context.Context) error {
	var err error

	for i := 0; i <= c.retry; i ++ {
		if err = c.client.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(c.retryBackoff)
	}

	return err
}
