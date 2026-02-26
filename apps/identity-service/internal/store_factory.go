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
	store, err := NewPostgresStore(databaseURL)
	if err != nil {
		return nil, err
	}
	if err := bootstrapAdminFromEnv(store); err != nil {
		return nil, err
	}
	return store, nil
}

func bootstrapAdminFromEnv(store *PostgresStore) error {
	email := os.Getenv("ADMIN_BOOTSTRAP_EMAIL")
	password := os.Getenv("ADMIN_BOOTSTRAP_PASSWORD")
	if email == "" || password == "" {
		return nil
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("hash bootstrap admin password: %w", err)
	}
	if err := store.UpsertAdminUser(email, passwordHash); err != nil {
		return err
	}
	return nil
}
