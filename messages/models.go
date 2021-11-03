package messages

import "context"

// Repo is responsible for all data transactions for messages
type Repo interface {
	UpsertMessage(context.Context, *Message) (*Message, error)
}

type MessageRequest struct {
	Ts         *string            `json:"ts"`
	Sender     *string            `json:"sender"`
	Msg        map[string]string `json:"message"`
	SentFromIP *string            `json:"sent-from-ip"`
	Priority   int               `json:"priority"`
}

type Message struct {
	Timestamp int
	Sender string
	Msg []byte
	SentFromIP string
	Priority int

}
