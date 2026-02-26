package runtimecfg

import "testing"

func TestPersistentStorageRequiredByMode(t *testing.T) {
	required := PersistentStorageRequiredFromEnv(func(key string) string {
		if key == "PERSISTENCE_MODE" {
			return "strict"
		}
		return ""
	})
	if !required {
		t.Fatalf("expected strict persistence mode to require persistent storage")
	}
}

func TestPersistentStorageRequiredByGoEnv(t *testing.T) {
	required := PersistentStorageRequiredFromEnv(func(key string) string {
		if key == "GO_ENV" {
			return "production"
		}
		return ""
	})
	if !required {
		t.Fatalf("expected production GO_ENV to require persistent storage")
	}
}

func TestPersistentStorageOptionalByDefault(t *testing.T) {
	required := PersistentStorageRequiredFromEnv(func(string) string {
		return ""
	})
	if required {
		t.Fatalf("expected persistent storage to be optional by default")
	}
}
