package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func (s *PostgresStore) ListWorkflows() ([]WorkflowTemplate, error) {
	rows, err := s.db.Query(
		`select id::text, name, description, content_suitability, age_band, steps, model_profile_id, safety_profile, version
		 from creator.workflow_templates
		 order by name asc`,
	)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	defer rows.Close()
	result := make([]WorkflowTemplate, 0, 32)
	for rows.Next() {
		var workflow WorkflowTemplate
		var rawSteps []byte
		if err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.ContentSuitability,
			&workflow.AgeBand,
			&rawSteps,
			&workflow.ModelProfileID,
			&workflow.SafetyProfile,
			&workflow.Version,
		); err != nil {
			return nil, fmt.Errorf("scan workflow: %w", err)
		}
		if err := json.Unmarshal(rawSteps, &workflow.Steps); err != nil {
			return nil, fmt.Errorf("decode workflow steps: %w", err)
		}
		result = append(result, workflow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workflows: %w", err)
	}
	return result, nil
}

func (s *PostgresStore) CreateWorkflow(workflow WorkflowTemplate, createdBy string) (WorkflowTemplate, error) {
	steps, err := json.Marshal(workflow.Steps)
	if err != nil {
		return WorkflowTemplate{}, fmt.Errorf("encode steps: %w", err)
	}
	var created WorkflowTemplate
	err = s.db.QueryRow(
		`insert into creator.workflow_templates
		 (name, description, content_suitability, age_band, steps, model_profile_id, safety_profile, version, created_by, updated_at)
		 values ($1, $2, $3, $4, $5::jsonb, $6, $7, 1, $8, now())
		 returning id::text, name, description, content_suitability, age_band, steps, model_profile_id, safety_profile, version`,
		workflow.Name,
		workflow.Description,
		workflow.ContentSuitability,
		workflow.AgeBand,
		steps,
		workflow.ModelProfileID,
		workflow.SafetyProfile,
		createdBy,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.ContentSuitability,
		&created.AgeBand,
		&steps,
		&created.ModelProfileID,
		&created.SafetyProfile,
		&created.Version,
	)
	if err != nil {
		return WorkflowTemplate{}, fmt.Errorf("insert workflow: %w", err)
	}
	if err := json.Unmarshal(steps, &created.Steps); err != nil {
		return WorkflowTemplate{}, fmt.Errorf("decode workflow steps: %w", err)
	}
	if err := s.snapshotWorkflowVersion(created, createdBy); err != nil {
		return WorkflowTemplate{}, err
	}
	return created, nil
}

func (s *PostgresStore) UpdateWorkflow(workflow WorkflowTemplate, updatedBy string) (WorkflowTemplate, bool, error) {
	steps, err := json.Marshal(workflow.Steps)
	if err != nil {
		return WorkflowTemplate{}, false, fmt.Errorf("encode steps: %w", err)
	}
	var updated WorkflowTemplate
	err = s.db.QueryRow(
		`update creator.workflow_templates
		 set name = $2,
		     description = $3,
		     content_suitability = $4,
		     age_band = $5,
		     steps = $6::jsonb,
		     model_profile_id = $7,
		     safety_profile = $8,
		     version = version + 1,
		     updated_at = now(),
		     created_by = coalesce(created_by, $9)
		 where id::text = $1
		 returning id::text, name, description, content_suitability, age_band, steps, model_profile_id, safety_profile, version`,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		workflow.ContentSuitability,
		workflow.AgeBand,
		steps,
		workflow.ModelProfileID,
		workflow.SafetyProfile,
		updatedBy,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.ContentSuitability,
		&updated.AgeBand,
		&steps,
		&updated.ModelProfileID,
		&updated.SafetyProfile,
		&updated.Version,
	)
	if err == sql.ErrNoRows {
		return WorkflowTemplate{}, false, nil
	}
	if err != nil {
		return WorkflowTemplate{}, false, fmt.Errorf("update workflow: %w", err)
	}
	if err := json.Unmarshal(steps, &updated.Steps); err != nil {
		return WorkflowTemplate{}, false, fmt.Errorf("decode workflow steps: %w", err)
	}
	if err := s.snapshotWorkflowVersion(updated, updatedBy); err != nil {
		return WorkflowTemplate{}, false, err
	}
	return updated, true, nil
}

func (s *PostgresStore) DeleteWorkflow(workflowID, deletedBy string) (bool, error) {
	result, err := s.db.Exec(`delete from creator.workflow_templates where id::text = $1`, workflowID)
	if err != nil {
		return false, fmt.Errorf("delete workflow: %w", err)
	}
	if err := s.writeAuditAction(deletedBy, "workflow_deleted", "workflow", workflowID, nil); err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("workflow rows affected: %w", err)
	}
	return affected > 0, nil
}

func (s *PostgresStore) FindWorkflow(workflowID string) (WorkflowTemplate, bool, error) {
	var workflow WorkflowTemplate
	var rawSteps []byte
	err := s.db.QueryRow(
		`select id::text, name, description, content_suitability, age_band, steps, model_profile_id, safety_profile, version
		 from creator.workflow_templates
		 where id::text = $1`,
		workflowID,
	).Scan(
		&workflow.ID,
		&workflow.Name,
		&workflow.Description,
		&workflow.ContentSuitability,
		&workflow.AgeBand,
		&rawSteps,
		&workflow.ModelProfileID,
		&workflow.SafetyProfile,
		&workflow.Version,
	)
	if err == sql.ErrNoRows {
		return WorkflowTemplate{}, false, nil
	}
	if err != nil {
		return WorkflowTemplate{}, false, fmt.Errorf("find workflow: %w", err)
	}
	if err := json.Unmarshal(rawSteps, &workflow.Steps); err != nil {
		return WorkflowTemplate{}, false, fmt.Errorf("decode workflow steps: %w", err)
	}
	return workflow, true, nil
}
