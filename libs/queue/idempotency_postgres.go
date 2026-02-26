package queue

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresIdempotencyStore struct {
	db            *sql.DB
	consumerScope string
	ttl           time.Duration
}

func newPostgresIdempotencyStore(
	databaseURL string,
	consumerScope string,
	ttl time.Duration,
) (*postgresIdempotencyStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &postgresIdempotencyStore{
		db:            db,
		consumerScope: consumerScope,
		ttl:           ttl,
	}, nil
}

func (s *postgresIdempotencyStore) Seen(eventID string) bool {
	if eventID == "" {
		return false
	}
	result, err := s.db.Exec(
		`insert into events.idempotency_keys (consumer_scope, event_id, expires_at)
		 values ($1, $2, $3)
		 on conflict (consumer_scope, event_id) do nothing`,
		s.consumerScope,
		eventID,
		time.Now().UTC().Add(s.ttl),
	)
	if err != nil {
		return false
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return affected == 0
}

func (s *postgresIdempotencyStore) Close() error {
	return s.db.Close()
}
