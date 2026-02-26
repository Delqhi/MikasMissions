package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const streamName = "MM_EVENTS"

type NATSBus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	mu   sync.Mutex
	subs []*nats.Subscription
}

func NewNATSBus(url string) (*NATSBus, error) {
	conn, err := nats.Connect(url, nats.Name("mikasmissions-platform"), nats.RetryOnFailedConnect(true))
	if err != nil {
		return nil, err
	}
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		conn.Close()
		return nil, err
	}
	bus := &NATSBus{conn: conn, js: js}
	if err := bus.ensureStream(); err != nil {
		_ = bus.Close()
		return nil, err
	}
	return bus, nil
}

func (b *NATSBus) ensureStream() error {
	if _, err := b.js.StreamInfo(streamName); err == nil {
		return nil
	}
	_, err := b.js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{">"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
	})
	if err != nil {
		if _, infoErr := b.js.StreamInfo(streamName); infoErr == nil {
			return nil
		}
	}
	return err
}

func (b *NATSBus) Publish(_ context.Context, event Event) error {
	msg := nats.NewMsg(event.Topic)
	msg.Data = event.Payload
	if event.ID != "" {
		msg.Header.Set("Nats-Msg-Id", event.ID)
	}
	_, err := b.js.PublishMsg(msg)
	return err
}

func (b *NATSBus) Subscribe(ctx context.Context, topic string, consumer string, handler Handler) error {
	sub, err := b.js.Subscribe(topic, func(msg *nats.Msg) {
		eventID := msg.Header.Get("Nats-Msg-Id")
		if eventID == "" {
			meta, metaErr := msg.Metadata()
			if metaErr == nil {
				eventID = fmt.Sprintf("%s-%d", topic, meta.Sequence.Stream)
			}
		}
		event := Event{ID: eventID, Topic: msg.Subject, Payload: msg.Data}
		if err := handler(context.Background(), event); err != nil {
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	},
		nats.Durable(consumer),
		nats.ManualAck(),
		nats.DeliverNew(),
		nats.AckExplicit(),
	)
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.subs = append(b.subs, sub)
	b.mu.Unlock()
	go func() {
		<-ctx.Done()
		_ = sub.Drain()
	}()
	return nil
}

func (b *NATSBus) Close() error {
	b.mu.Lock()
	for _, sub := range b.subs {
		_ = sub.Drain()
	}
	b.subs = nil
	b.mu.Unlock()
	b.conn.Drain()
	b.conn.Close()
	return nil
}
