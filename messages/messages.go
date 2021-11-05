package messages

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/unity-web-service/queue"
	"github.com/unity-web-service/store/postgres"
)

type MessageOps struct {
	pg *postgres.Conn
}

func New(db *postgres.Conn) *MessageOps {
	return &MessageOps{pg: db}
}

func (m *MessageOps) UpsertMessage(ctx context.Context, data *Message) error {
	tx, err := m.pg.Client.Begin()
	if err != nil {
		return err
	}

	if _, err := m.pg.Queries.ExecContext(ctx, tx, "upsert-message",
		data.Timestamp,
		data.Priority,
		data.Sender,
		data.SentFromIP,
		data.Msg,
		data.Queued,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MessageOps) GetAllUnqueuedMessages(ctx context.Context) ([]*Message, error) {
	rows, err := m.pg.Queries.QueryContext(ctx, m.pg.Client, "get-all-unqueued-messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uqMsgs := make([]*Message, 0)
	for rows.Next() {
		msg := new(Message)
		if err := rows.Scan(&msg.Timestamp, &msg.Priority, &msg.Sender, &msg.SentFromIP, &msg.Msg, &msg.Queued); err != nil {
			return nil, err
		}
		uqMsgs = append(uqMsgs, msg)
	}

	return uqMsgs, rows.Err()
}

// Reconcile gets all unqueued messages and attemps to queue them and update the queued status in postgres
func (m *MessageOps) Reconcile(ctx context.Context, publisher queue.IPublisher) {
	for range time.Tick(15 * time.Second) {
		msgs, err := m.GetAllUnqueuedMessages(ctx)
		if err != nil {
			log.Printf("could not retrieved unqueued messages: %v", err)
		}
		log.Printf("messages: %v", msgs)
		for _, msg := range msgs {
			// we don't want to
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("could not marshal message struct: %v", err)
				continue
			}
			// we don't want to update a message in postgres that hasn't published, so continue
			if err := publisher.Publish(b); err != nil {
				log.Printf("could not publish message: %v", err)
				continue
			}
			// update message with queued=true
			msg.Queued = true
			// we don't want to leave all the other messages hanging if this error is a retryable one
			if err := m.UpsertMessage(ctx, msg); err != nil {
				log.Printf("could not upsert message: %v", err)
				continue
			}
		}
	}
}
