package internal

import "encoding/json"

type WorkflowRun struct {
	ID           string
	WorkflowID   string
	Status       string
	Priority     string
	AutoPublish  bool
	InputPayload json.RawMessage
	LastError    string
}

type WorkflowRunLog struct {
	RunID     string
	Step      string
	Status    string
	Message   string
	EventTime string
}
