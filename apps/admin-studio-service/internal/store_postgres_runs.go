package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func (s *PostgresStore) CreateRun(run WorkflowRun, createdBy string) (WorkflowRun, error) {
	if len(run.InputPayload) == 0 {
		run.InputPayload = json.RawMessage(`{}`)
	}
	var created WorkflowRun
	err := s.db.QueryRow(
		`insert into creator.workflow_runs
		 (workflow_id, status, input_payload, priority, auto_publish, created_by, updated_at)
		 values ($1::uuid, 'requested', $2::jsonb, $3, $4, $5, now())
		 returning id::text, workflow_id::text, status, priority, auto_publish, input_payload, coalesce(last_error, '')`,
		run.WorkflowID,
		run.InputPayload,
		run.Priority,
		run.AutoPublish,
		createdBy,
	).Scan(
		&created.ID,
		&created.WorkflowID,
		&created.Status,
		&created.Priority,
		&created.AutoPublish,
		&created.InputPayload,
		&created.LastError,
	)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("insert workflow run: %w", err)
	}
	return created, nil
}

func (s *PostgresStore) FindRun(runID string) (WorkflowRun, bool, error) {
	var run WorkflowRun
	err := s.db.QueryRow(
		`select id::text, workflow_id::text, status, priority, auto_publish, input_payload, coalesce(last_error, '')
		 from creator.workflow_runs
		 where id::text = $1`,
		runID,
	).Scan(
		&run.ID,
		&run.WorkflowID,
		&run.Status,
		&run.Priority,
		&run.AutoPublish,
		&run.InputPayload,
		&run.LastError,
	)
	if err == sql.ErrNoRows {
		return WorkflowRun{}, false, nil
	}
	if err != nil {
		return WorkflowRun{}, false, fmt.Errorf("find run: %w", err)
	}
	return run, true, nil
}

func (s *PostgresStore) ListRunLogs(runID string) ([]WorkflowRunLog, error) {
	rows, err := s.db.Query(
		`select run_id::text, step, status, message, created_at
		 from creator.workflow_run_steps
		 where run_id::text = $1
		 order by id asc`,
		runID,
	)
	if err != nil {
		return nil, fmt.Errorf("list run logs: %w", err)
	}
	defer rows.Close()
	logs := make([]WorkflowRunLog, 0, 16)
	for rows.Next() {
		var item WorkflowRunLog
		if err := rows.Scan(&item.RunID, &item.Step, &item.Status, &item.Message, &item.EventTime); err != nil {
			return nil, fmt.Errorf("scan run log: %w", err)
		}
		logs = append(logs, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate run logs: %w", err)
	}
	return logs, nil
}

func (s *PostgresStore) AppendRunLog(logEntry WorkflowRunLog) error {
	_, err := s.db.Exec(
		`insert into creator.workflow_run_steps (run_id, step, status, message)
		 values ($1::uuid, $2, $3, $4)`,
		logEntry.RunID,
		logEntry.Step,
		logEntry.Status,
		logEntry.Message,
	)
	if err != nil {
		return fmt.Errorf("insert run log: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetRunStatus(runID, status, lastError string) (bool, error) {
	result, err := s.db.Exec(
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
		return false, fmt.Errorf("update run status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("run status rows affected: %w", err)
	}
	return affected > 0, nil
}
