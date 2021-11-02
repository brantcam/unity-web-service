package messages

import "context"

// Repo is responsible for all data transactions for messages
type Repo interface {
	UpsertMessage(context.Context, *Message) (*Message, error)
}

type Message struct {
	Ts         string            `json:"ts,omitempty"`
	Sender     string            `json:"sender,omitempty"`
	Msg        map[string]string `json:"message,omitempty"`
	SentFromIP string            `json:"sent-from-ip,omitempty"`
	Priority   int               `json:"priority,omitempty"`
}
