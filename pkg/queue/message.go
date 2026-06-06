package queue

import "context"

type Event interface {
	Subject() string
}

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(subject string, handler func(ctx context.Context, event []byte) error) error
	Close() error
}
