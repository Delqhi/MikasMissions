package queue

import "context"

type Event struct {
	ID      string
	Topic   string
	Payload []byte
}

type Handler func(ctx context.Context, event Event) error

type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, topic string, consumer string, handler Handler) error
	Close() error
}
