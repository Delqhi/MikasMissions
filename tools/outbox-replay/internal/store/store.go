package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type FailedEvent struct {
	EventID     string
	Topic       string
	Attempts    int
	LastError   string
	AvailableAt time.Time
	CreatedAt   time.Time
}

type Store struct {
	db *sql.DB
}

func New(databaseURL string) (*Store, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ListFailed(ctx context.Context, topic string, limit int) ([]FailedEvent, error) {
	query := `select event_id, topic, attempts, coalesce(last_error, ''), available_at, created_at
		from events.outbox
		where status = 'failed'`
	args := make([]any, 0, 2)
	if topic != "" {
		query += " and topic = $1"
		args = append(args, topic)
	}
	query += fmt.Sprintf(" order by id desc limit $%d", len(args)+1)
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed outbox rows: %w", err)
	}
	defer rows.Close()
	result := make([]FailedEvent, 0, limit)
	for rows.Next() {
		var row FailedEvent
		if err := rows.Scan(&row.EventID, &row.Topic, &row.Attempts, &row.LastError, &row.AvailableAt, &row.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan failed outbox row: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate failed outbox rows: %w", err)
	}
	return result, nil
}

func (s *Store) RequeueFailed(
	ctx context.Context,
	topic string,
	limit int,
	dryRun bool,
	resetAttempts bool,
) ([]string, error) {
	ids, err := s.selectFailedEventIDs(ctx, topic, limit)
	if err != nil {
		return nil, err
	}
	if dryRun || len(ids) == 0 {
		return ids, nil
	}
	return s.requeueEventIDs(ctx, ids, resetAttempts)
}

func (s *Store) RequeueEvent(ctx context.Context, eventID string, dryRun bool, resetAttempts bool) ([]string, error) {
	if strings.TrimSpace(eventID) == "" {
		return nil, fmt.Errorf("eventID is required")
	}
	if dryRun {
		exists, err := s.hasFailedRow(ctx, eventID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return []string{}, nil
		}
		return []string{eventID}, nil
	}
	return s.requeueEventIDs(ctx, []string{eventID}, resetAttempts)
}

func (s *Store) selectFailedEventIDs(ctx context.Context, topic string, limit int) ([]string, error) {
	query := `select event_id from events.outbox where status = 'failed'`
	args := make([]any, 0, 2)
	if topic != "" {
		query += " and topic = $1"
		args = append(args, topic)
	}
	query += fmt.Sprintf(" order by id asc limit $%d", len(args)+1)
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed requeue candidates: %w", err)
	}
	defer rows.Close()
	ids := make([]string, 0, limit)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan failed requeue candidate: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate failed requeue candidates: %w", err)
	}
	return ids, nil
}

func (s *Store) requeueEventIDs(ctx context.Context, eventIDs []string, resetAttempts bool) ([]string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin requeue tx: %w", err)
	}
	defer tx.Rollback()
	requeued := make([]string, 0, len(eventIDs))
	for _, eventID := range eventIDs {
		result, err := tx.ExecContext(
			ctx,
			`update events.outbox
			 set status = 'pending',
			     available_at = now(),
			     last_error = null,
			     attempts = case when $2 then 0 else attempts end
			 where event_id = $1 and status = 'failed'`,
			eventID,
			resetAttempts,
		)
		if err != nil {
			return nil, fmt.Errorf("requeue outbox row %s: %w", eventID, err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("read requeue rows affected for %s: %w", eventID, err)
		}
		if affected > 0 {
			requeued = append(requeued, eventID)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit requeue tx: %w", err)
	}
	return requeued, nil
}

func (s *Store) hasFailedRow(ctx context.Context, eventID string) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(
		ctx,
		`select exists(
			select 1 from events.outbox where event_id = $1 and status = 'failed'
		)`,
		eventID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("query failed row existence: %w", err)
	}
	return exists, nil
}
