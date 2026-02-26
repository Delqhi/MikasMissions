package queue

import "context"

type OutboxItem struct {
	Event Event
	Sent  bool
}

type Outbox struct {
	items []OutboxItem
}

func NewOutbox() *Outbox {
	return &Outbox{items: []OutboxItem{}}
}

func (o *Outbox) Add(event Event) {
	o.items = append(o.items, OutboxItem{Event: event})
}

func (o *Outbox) Flush(ctx context.Context, bus Bus) error {
	for i := range o.items {
		if o.items[i].Sent {
			continue
		}
		if err := bus.Publish(ctx, o.items[i].Event); err != nil {
			return err
		}
		o.items[i].Sent = true
	}
	return nil
}
