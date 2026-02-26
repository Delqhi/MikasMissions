package queue

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PersistentOutbox struct {
	db *sql.DB
}

func NewPersistentOutbox(databaseURL string) (*PersistentOutbox, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &PersistentOutbox{db: db}, nil
}

func (o *PersistentOutbox) Add(event Event) error {
	if event.ID == "" {
		return fmt.Errorf("event id is required for persistent outbox")
	}
	_, err := o.db.Exec(
		`insert into events.outbox (event_id, topic, payload, status, available_at)
		 values ($1, $2, $3, 'pending', now())
		 on conflict (event_id) do nothing`,
		event.ID, event.Topic, event.Payload,
	)
	if err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}
	return nil
}

func (o *PersistentOutbox) Close() error {
	return o.db.Close()
}
