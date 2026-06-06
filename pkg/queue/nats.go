package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSConfig struct {
	URL string
}

type NATSBus struct {
	conn *nats.Conn
}

func NewNATS(cfg NATSConfig) (*NATSBus, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	return &NATSBus{conn: conn}, nil
}

func (b *NATSBus) Publish(ctx context.Context, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return b.conn.Publish(event.Subject(), data)
}

func (b *NATSBus) Subscribe(subject string, handler func(ctx context.Context, data []byte) error) error {
	_, err := b.conn.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(context.Background(), msg.Data); err != nil {
			return
		}
	})
	return err
}

func (b *NATSBus) Close() error {
	b.conn.Close()
	return nil
}
