package contractsevents

import (
	"encoding/json"
	"errors"
)

type VideoRunRequestedV1 struct {
	RunID              string          `json:"run_id"`
	WorkflowID         string          `json:"workflow_id"`
	ModelProfileID     string          `json:"model_profile_id"`
	InputPayload       json.RawMessage `json:"input_payload"`
	AutoPublish        bool            `json:"auto_publish"`
	Priority           string          `json:"priority"`
	ContentSuitability string          `json:"content_suitability"`
	AgeBand            string          `json:"age_band"`
	RequestedBy        string          `json:"requested_by"`
	RequestedAt        string          `json:"requested_at"`
	TraceID            string          `json:"trace_id"`
}

func (e VideoRunRequestedV1) Validate() error {
	if e.RunID == "" || e.WorkflowID == "" || e.ModelProfileID == "" || e.RequestedBy == "" {
		return errors.New("video.run.requested.v1 has missing required fields")
	}
	if e.RequestedAt == "" || e.TraceID == "" || e.ContentSuitability == "" || e.AgeBand == "" {
		return errors.New("video.run.requested.v1 has missing metadata fields")
	}
	if e.Priority == "" {
		return errors.New("video.run.requested.v1 requires priority")
	}
	return nil
}
