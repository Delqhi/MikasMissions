package internal

import (
	"context"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

type memoryOutboxWriter struct{}

func (w *memoryOutboxWriter) EnqueueAndFlush(ctx context.Context, bus queue.Bus, event queue.Event) error {
	outbox := queue.NewOutbox()
	outbox.Add(event)
	return outbox.Flush(ctx, bus)
}
