package internal

import "encoding/json"

type Repository interface {
	ListWorkflows() ([]WorkflowTemplate, error)
	CreateWorkflow(workflow WorkflowTemplate, createdBy string) (WorkflowTemplate, error)
	UpdateWorkflow(workflow WorkflowTemplate, updatedBy string) (WorkflowTemplate, bool, error)
	DeleteWorkflow(workflowID, deletedBy string) (bool, error)
	FindWorkflow(workflowID string) (WorkflowTemplate, bool, error)
	CreateRun(run WorkflowRun, createdBy string) (WorkflowRun, error)
	FindRun(runID string) (WorkflowRun, bool, error)
	ListRunLogs(runID string) ([]WorkflowRunLog, error)
	AppendRunLog(log WorkflowRunLog) error
	SetRunStatus(runID, status, lastError string) (bool, error)
	GetModelProfile(modelProfileID string) (ModelProfile, bool, error)
	PutModelProfile(profile ModelProfile, updatedBy string) (ModelProfile, error)
}

type runRequestedPayload struct {
	RunID              string          `json:"run_id"`
	WorkflowID         string          `json:"workflow_id"`
	ModelProfileID     string          `json:"model_profile_id"`
	InputPayload       json.RawMessage `json:"input_payload"`
	AutoPublish        bool            `json:"auto_publish"`
	Priority           string          `json:"priority"`
	ContentSuitability string          `json:"content_suitability"`
	AgeBand            string          `json:"age_band"`
	RequestedBy        string          `json:"requested_by"`
}
