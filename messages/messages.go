package messages

import (
	"context"

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

	if _, err := m.pg.Queries.ExecContext(ctx, tx, "insert-message",
		data.Timestamp,
		data.Priority,
		data.Sender,
		data.SentFromIP,
		data.Msg,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MessageOps) UpdateMessage(ctx context.Context, ts int, sender string) error {
	return nil
}

func (m *MessageOps) GetMessage(ctx context.Context, ts int, sender string) (*Message, error) {
	return nil, nil
}

func (m *MessageOps) GetUnqueuedMessages(ctx context.Context) ([]*Message, error) {
	return nil, nil
}
