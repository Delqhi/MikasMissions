package internal

import "testing"

func TestOpenRepositoryFromEnvDefaultsToInMemory(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	repo, err := OpenRepositoryFromEnv()
	if err != nil {
		t.Fatalf("open repository: %v", err)
	}
	if _, ok := repo.(*Store); !ok {
		t.Fatalf("expected in-memory store by default")
	}
}

func TestOpenRepositoryFromEnvStrictModeRequiresDatabaseURL(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("DATABASE_URL", "")
	repo, err := OpenRepositoryFromEnv()
	if err == nil {
		t.Fatalf("expected error, got repository %T", repo)
	}
}
