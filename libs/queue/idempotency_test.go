package queue

import "testing"

func TestIdempotencyGuardSeenOnce(t *testing.T) {
	guard := NewIdempotencyGuard()
	if guard.Seen("evt-1") {
		t.Fatalf("first occurrence should not be marked as duplicate")
	}
	if !guard.Seen("evt-1") {
		t.Fatalf("second occurrence should be marked as duplicate")
	}
}

func TestScopedIdempotencyGuardFallsBackWhenDatabaseUnavailable(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://127.0.0.1:1/does_not_exist?sslmode=disable")
	guard := NewScopedIdempotencyGuard("worker-test")
	if guard.Seen("evt-2") {
		t.Fatalf("first occurrence should not be marked as duplicate")
	}
	if !guard.Seen("evt-2") {
		t.Fatalf("second occurrence should be marked as duplicate")
	}
}

func TestScopedIdempotencyGuardPanicsInStrictModeWithoutDatabaseURL(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("DATABASE_URL", "")
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("expected strict mode panic")
		}
	}()
	_ = NewScopedIdempotencyGuard("worker-test")
}
