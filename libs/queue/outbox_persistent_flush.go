package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	outboxBatchSize   = 100
	outboxMaxAttempts = 10
	outboxRetryDelay  = 30 * time.Second
)

type pendingOutboxRow struct {
	ID       int64
	EventID  string
	Topic    string
	Payload  []byte
	Attempts int
}

func (o *PersistentOutbox) Flush(ctx context.Context, bus Bus) error {
	pending, err := o.pendingEvents(ctx)
	if err != nil {
		return err
	}
	for _, row := range pending {
		if publishErr := bus.Publish(ctx, Event{
			ID:      row.EventID,
			Topic:   row.Topic,
			Payload: row.Payload,
		}); publishErr != nil {
			if row.Attempts+1 >= outboxMaxAttempts {
				if err := o.markTerminalFailed(ctx, row.ID, publishErr); err != nil {
					return err
				}
				if err := o.publishDLQ(ctx, bus, row, publishErr); err != nil {
					return err
				}
				continue
			}
			if err := o.markRetryPending(ctx, row.ID, publishErr); err != nil {
				return err
			}
			continue
		}
		if err := o.markPublished(ctx, row.ID); err != nil {
			return err
		}
	}
	return nil
}

func (o *PersistentOutbox) pendingEvents(ctx context.Context) ([]pendingOutboxRow, error) {
	rows, err := o.db.QueryContext(
		ctx,
		`select id, event_id, topic, payload, attempts
		 from events.outbox
		 where status = 'pending'
		   and available_at <= now()
		 order by id asc
		 limit $1`,
		outboxBatchSize,
	)
	if err != nil {
		return nil, fmt.Errorf("query pending outbox events: %w", err)
	}
	defer rows.Close()
	result := make([]pendingOutboxRow, 0, 32)
	for rows.Next() {
		var row pendingOutboxRow
		if err := rows.Scan(&row.ID, &row.EventID, &row.Topic, &row.Payload, &row.Attempts); err != nil {
			return nil, fmt.Errorf("scan pending outbox event: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending outbox events: %w", err)
	}
	return result, nil
}

func (o *PersistentOutbox) markPublished(ctx context.Context, rowID int64) error {
	_, err := o.db.ExecContext(
		ctx,
		`update events.outbox
		 set status = 'published',
		     attempts = attempts + 1,
		     last_error = null,
		     published_at = now()
		 where id = $1`,
		rowID,
	)
	if err != nil {
		return fmt.Errorf("mark outbox published: %w", err)
	}
	return nil
}

func (o *PersistentOutbox) markRetryPending(ctx context.Context, rowID int64, publishErr error) error {
	_, err := o.db.ExecContext(
		ctx,
		`update events.outbox
		 set status = 'pending',
		     attempts = attempts + 1,
		     last_error = $2,
		     available_at = $3
		 where id = $1`,
		rowID,
		publishErr.Error(),
		time.Now().UTC().Add(outboxRetryDelay),
	)
	if err != nil {
		return fmt.Errorf("mark outbox retry pending: %w", err)
	}
	return nil
}

func (o *PersistentOutbox) markTerminalFailed(ctx context.Context, rowID int64, publishErr error) error {
	_, err := o.db.ExecContext(
		ctx,
		`update events.outbox
		 set status = 'failed',
		     attempts = attempts + 1,
		     last_error = $2
		 where id = $1`,
		rowID,
		publishErr.Error(),
	)
	if err != nil {
		return fmt.Errorf("mark outbox terminal failed: %w", err)
	}
	return nil
}

type outboxDLQEvent struct {
	EventID       string          `json:"event_id"`
	OriginalTopic string          `json:"original_topic"`
	Error         string          `json:"error"`
	Attempts      int             `json:"attempts"`
	FailedAt      string          `json:"failed_at"`
	Payload       json.RawMessage `json:"payload"`
}

func (o *PersistentOutbox) publishDLQ(ctx context.Context, bus Bus, row pendingOutboxRow, publishErr error) error {
	payload, err := json.Marshal(outboxDLQEvent{
		EventID:       row.EventID,
		OriginalTopic: row.Topic,
		Error:         publishErr.Error(),
		Attempts:      row.Attempts + 1,
		FailedAt:      time.Now().UTC().Format(time.RFC3339),
		Payload:       json.RawMessage(row.Payload),
	})
	if err != nil {
		return fmt.Errorf("marshal outbox dlq payload: %w", err)
	}
	if err := bus.Publish(ctx, Event{
		ID:      row.EventID + "-dlq",
		Topic:   row.Topic + ".dlq.v1",
		Payload: payload,
	}); err != nil {
		return fmt.Errorf("publish outbox dlq event: %w", err)
	}
	return nil
}
