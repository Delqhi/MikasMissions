package internal

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/delqhi/mikasmissions/platform/libs/runtimecfg"
)

type outboxWriter interface {
	EnqueueAndFlush(ctx context.Context, bus queue.Bus, event queue.Event) error
}

func newOutboxWriterFromEnv() (outboxWriter, io.Closer, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		if runtimecfg.PersistentStorageRequired() {
			return nil, nil, fmt.Errorf("DATABASE_URL is required when persistent storage is strict")
		}
		return &memoryOutboxWriter{}, nil, nil
	}
	persistent, err := queue.NewPersistentOutbox(databaseURL)
	if err != nil {
		return nil, nil, err
	}
	return &persistentOutboxWriter{outbox: persistent}, persistent, nil
}
