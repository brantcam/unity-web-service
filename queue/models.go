package queue

type IPublisher interface {
	Publish(m []byte) error
}