package messages

import "context"

// Repo is responsible for all data transactions for messages
type Repo interface {
	UpsertMessage(context.Context, *Message) error
	GetAllUnqueuedMessages(ctx context.Context) ([]*Message, error)
}

type MessageRequest struct {
	Ts         *string           `json:"ts"`
	Sender     *string           `json:"sender"`
	SentFromIP *string           `json:"sent-from-ip"`
	Msg        map[string]string `json:"message"`
	Priority   int               `json:"priority"`
}

type Message struct {
	Timestamp  int
	Priority   int
	Sender     string
	SentFromIP string
	Msg        string
	Queued     bool
}
