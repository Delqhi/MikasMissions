package internal

import (
	"encoding/json"
	"fmt"
)

func (s *PostgresStore) snapshotWorkflowVersion(workflow WorkflowTemplate, actor string) error {
	snapshot, err := json.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("encode workflow snapshot: %w", err)
	}
	_, err = s.db.Exec(
		`insert into creator.workflow_template_versions
		 (workflow_id, version, snapshot, created_by)
		 values ($1::uuid, $2, $3::jsonb, $4)`,
		workflow.ID,
		workflow.Version,
		snapshot,
		actor,
	)
	if err != nil {
		return fmt.Errorf("insert workflow version snapshot: %w", err)
	}
	return nil
}

func (s *PostgresStore) writeAuditAction(actorID, action, resourceType, resourceID string, payload []byte) error {
	if actorID == "" {
		actorID = "system"
	}
	if payload == nil {
		payload = []byte(`{}`)
	}
	_, err := s.db.Exec(
		`insert into audit.admin_actions
		 (admin_user_id, action, resource_type, resource_id, payload)
		 values ($1, $2, $3, $4, $5::jsonb)`,
		actorID,
		action,
		resourceType,
		resourceID,
		payload,
	)
	if err != nil {
		return fmt.Errorf("insert admin audit action: %w", err)
	}
	return nil
}
