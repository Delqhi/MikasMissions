package contractsevents

import "errors"

type VideoWorkflowCreatedV1 struct {
	WorkflowID string `json:"workflow_id"`
	Version    int    `json:"version"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

func (e VideoWorkflowCreatedV1) Validate() error {
	if e.WorkflowID == "" || e.CreatedBy == "" || e.CreatedAt == "" {
		return errors.New("video.workflow.created.v1 has missing required fields")
	}
	if e.Version <= 0 {
		return errors.New("video.workflow.created.v1 requires positive version")
	}
	return nil
}
