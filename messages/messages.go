package messages

import (
	"context"
	"time"

	"github.com/unity-web-service/store/postgres"
)

type MessageOps struct {
	pg *postgres.Conn
}

func New(db *postgres.Conn) *MessageOps {
	return &MessageOps{pg: db}
}

func (m *MessageOps) InsertMessage(ctx context.Context, data *Message) error {
	tx, err := m.pg.Client.Begin()
	if err != nil {
		return err
	}
	for i := 0; i <= m.pg.Retry; i++ {
		if _, err = m.pg.Queries.ExecContext(ctx, tx, "insert-message",
			data.Timestamp,
			data.Priority,
			data.Sender,
			data.SentFromIP,
			data.Msg,
		); err == nil {
			break
		}
		time.Sleep(m.pg.RetryBackoff)
	}
	// checking the err of the ExecContext
	if err != nil {
		return err
	}

	return tx.Commit()
}
