package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/dotsql"
	"github.com/unity-web-service/config"

	// postgres driver for the database/sql package
	_ "github.com/lib/pq"
)

type Conn struct {
	Queries *dotsql.DotSql
	Client  *sqlx.DB

	Retry        int
	RetryBackoff time.Duration
}

func New(ctx context.Context, c config.Config) (*Conn, error) {
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
		Queries:      dot,
		Client:       db,
		Retry:        3,
		RetryBackoff: 5 * time.Second,
	}

	return conn, db.PingContext(ctx)
}

func (c *Conn) Health(ctx context.Context) error {
	var err error

	for i := 0; i <= c.Retry; i++ {
		if err = c.Client.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(c.RetryBackoff)
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) MigrateUp(ctx context.Context) error {
	tx, err := c.Client.Begin()
	if err != nil {
		return err
	}

	if _, err := c.Queries.ExecContext(ctx, tx, "create-messsages-table"); err != nil {
		return err
	}

	return tx.Commit()
}
