package internal

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

type stubOutbox struct {
	calls int
	err   error
}

func (s *stubOutbox) Flush(_ context.Context, _ queue.Bus) error {
	s.calls++
	return s.err
}

func TestFlushOnceDelegatesToOutbox(t *testing.T) {
	outbox := &stubOutbox{}
	relay := NewRelay(outbox, queue.NewInMemoryBus(), 10*time.Millisecond, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	if err := relay.FlushOnce(context.Background()); err != nil {
		t.Fatalf("flush once: %v", err)
	}
	if outbox.calls != 1 {
		t.Fatalf("expected one outbox call, got %d", outbox.calls)
	}
}

func TestFlushOnceReturnsOutboxError(t *testing.T) {
	expectedErr := errors.New("flush failed")
	outbox := &stubOutbox{err: expectedErr}
	relay := NewRelay(outbox, queue.NewInMemoryBus(), 10*time.Millisecond, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	if err := relay.FlushOnce(context.Background()); !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestRunCallsFlushRepeatedly(t *testing.T) {
	outbox := &stubOutbox{}
	relay := NewRelay(outbox, queue.NewInMemoryBus(), 1*time.Millisecond, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Millisecond)
	defer cancel()
	relay.Run(ctx)
	if outbox.calls < 2 {
		t.Fatalf("expected at least 2 outbox flush calls, got %d", outbox.calls)
	}
}
