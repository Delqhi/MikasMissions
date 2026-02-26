package config

import (
	"strings"
	"testing"
	"time"
)

func TestParseDefaultsAndEnvDatabaseURL(t *testing.T) {
	opts, err := Parse([]string{}, func(key string) string {
		if key == "DATABASE_URL" {
			return "postgres://example"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if opts.DatabaseURL != "postgres://example" {
		t.Fatalf("unexpected database url: %s", opts.DatabaseURL)
	}
	if opts.Mode != ModeListFailed {
		t.Fatalf("unexpected default mode: %s", opts.Mode)
	}
	if opts.DryRun != true {
		t.Fatalf("expected dry-run true by default")
	}
	if opts.ResetAttempts != true {
		t.Fatalf("expected reset-attempts true by default")
	}
	if opts.Timeout != 15*time.Second {
		t.Fatalf("unexpected default timeout: %s", opts.Timeout)
	}
}

func TestParseRejectsUnsupportedMode(t *testing.T) {
	_, err := Parse([]string{"-database-url", "postgres://example", "-mode", "invalid"}, func(string) string {
		return ""
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported mode") {
		t.Fatalf("expected unsupported mode error, got: %v", err)
	}
}

func TestParseRejectsMissingEventIDInEventMode(t *testing.T) {
	_, err := Parse([]string{"-database-url", "postgres://example", "-mode", "requeue-event"}, func(string) string {
		return ""
	})
	if err == nil || !strings.Contains(err.Error(), "event-id is required") {
		t.Fatalf("expected event-id required error, got: %v", err)
	}
}

func TestParseRejectsMissingDatabaseURL(t *testing.T) {
	_, err := Parse([]string{}, func(string) string {
		return ""
	})
	if err == nil || !strings.Contains(err.Error(), "database URL is required") {
		t.Fatalf("expected database URL required error, got: %v", err)
	}
}
