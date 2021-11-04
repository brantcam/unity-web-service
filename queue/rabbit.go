package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Queue string
}

func (p *Publisher) Publish(m []byte) error {
	conn, err := amqp.Dial(fmt.Sprintf(`amqp://%s:%s@%s:%s/`, p.User, p.Pass, p.Host, p.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(p.Queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	return ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        m,
	})
}
