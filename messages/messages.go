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

func (m *MessageOps) InsertMessage(ctx context.Context, data *Message) (*Message, error) {
	return nil, m.pg.Health(context.Background())
}
