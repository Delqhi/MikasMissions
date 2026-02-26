package internal

import (
	"fmt"
	"os"

	"github.com/delqhi/mikasmissions/platform/libs/runtimecfg"
)

func OpenRepositoryFromEnv() (Repository, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		if runtimecfg.PersistentStorageRequired() {
			return nil, fmt.Errorf("DATABASE_URL is required when persistent storage is strict")
		}
		return NewStore(), nil
	}
	return NewPostgresStore(databaseURL)
}
