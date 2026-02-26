package internal

import "testing"

func TestNewOutboxWriterFromEnvDefaultsToMemory(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	writer, closer, err := newOutboxWriterFromEnv()
	if err != nil {
		t.Fatalf("new outbox writer: %v", err)
	}
	if writer == nil {
		t.Fatalf("expected outbox writer")
	}
	if closer != nil {
		t.Fatalf("expected no closer for memory writer")
	}
}

func TestNewOutboxWriterFromEnvStrictModeRequiresDatabaseURL(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("DATABASE_URL", "")
	writer, closer, err := newOutboxWriterFromEnv()
	if err == nil {
		t.Fatalf("expected strict mode error, got writer=%T closer=%T", writer, closer)
	}
}
