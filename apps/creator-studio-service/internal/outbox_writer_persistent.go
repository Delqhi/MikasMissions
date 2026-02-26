package internal

import (
	"context"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

type persistentOutboxWriter struct {
	outbox *queue.PersistentOutbox
}

func (w *persistentOutboxWriter) EnqueueAndFlush(_ context.Context, _ queue.Bus, event queue.Event) error {
	return w.outbox.Add(event)
}
