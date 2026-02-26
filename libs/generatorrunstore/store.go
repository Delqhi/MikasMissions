package generatorrunstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	openOnce sync.Once
	openErr  error
	db       *sql.DB
)

func SetRunStatus(ctx context.Context, runID, status, lastError string) error {
	handle, err := openDB()
	if err != nil || handle == nil {
		return err
	}
	_, err = handle.ExecContext(
		ctx,
		`update creator.workflow_runs
		 set status = $2,
		     last_error = $3,
		     updated_at = now()
		 where id::text = $1`,
		runID,
		status,
		lastError,
	)
	if err != nil {
		return fmt.Errorf("update workflow run status: %w", err)
	}
	return nil
}

func AppendRunLog(ctx context.Context, runID, step, status, message string) error {
	handle, err := openDB()
	if err != nil || handle == nil {
		return err
	}
	_, err = handle.ExecContext(
		ctx,
		`insert into creator.workflow_run_steps (run_id, step, status, message)
		 values ($1::uuid, $2, $3, $4)`,
		runID,
		step,
		status,
		message,
	)
	if err != nil {
		return fmt.Errorf("insert run log: %w", err)
	}
	return nil
}

func openDB() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil
	}
	openOnce.Do(func() {
		handle, err := sql.Open("pgx", databaseURL)
		if err != nil {
			openErr = fmt.Errorf("open postgres: %w", err)
			return
		}
		if err := handle.Ping(); err != nil {
			_ = handle.Close()
			openErr = fmt.Errorf("ping postgres: %w", err)
			return
		}
		db = handle
	})
	return db, openErr
}
