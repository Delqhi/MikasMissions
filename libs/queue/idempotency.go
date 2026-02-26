package queue

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/runtimecfg"
)

type idempotencyStore interface {
	Seen(eventID string) bool
}

type IdempotencyGuard struct {
	store idempotencyStore
}

func NewIdempotencyGuard() *IdempotencyGuard {
	return &IdempotencyGuard{store: newMemoryIdempotencyStore()}
}

func NewScopedIdempotencyGuard(consumerScope string) *IdempotencyGuard {
	if consumerScope == "" {
		return NewIdempotencyGuard()
	}
	strict := runtimecfg.PersistentStorageRequired()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		if strict {
			panic("DATABASE_URL is required for scoped idempotency guard in strict persistence mode")
		}
		return NewIdempotencyGuard()
	}
	store, err := newPostgresIdempotencyStore(databaseURL, consumerScope, 7*24*time.Hour)
	if err != nil {
		if strict {
			panic(fmt.Sprintf("open scoped idempotency store: %v", err))
		}
		return NewIdempotencyGuard()
	}
	return &IdempotencyGuard{store: store}
}

func (g *IdempotencyGuard) Seen(eventID string) bool {
	return g.store.Seen(eventID)
}

func (g *IdempotencyGuard) Close() error {
	closer, ok := g.store.(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}

type memoryIdempotencyStore struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func newMemoryIdempotencyStore() *memoryIdempotencyStore {
	return &memoryIdempotencyStore{seen: make(map[string]struct{})}
}

func (s *memoryIdempotencyStore) Seen(eventID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seen[eventID]; ok {
		return true
	}
	s.seen[eventID] = struct{}{}
	return false
}
