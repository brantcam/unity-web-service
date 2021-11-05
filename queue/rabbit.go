package queue

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Queue string

	Retry        int
	RetryBackoff time.Duration
}

// Publish will send messages to the Queue and retry if an error occurs
func (p *Publisher) Publish(m []byte) error {
	conn, err := amqp.Dial(fmt.Sprintf(
		`amqp://%s:%s@%s:%s/`,
		p.User,
		p.Pass,
		p.Host,
		p.Port),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(p.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	for i := 0; i <= p.Retry; i++ {
		if err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        m,
		}); err == nil {
			break
		}
		time.Sleep(p.RetryBackoff)
	}

	return err
}
