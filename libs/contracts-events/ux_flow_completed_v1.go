package contractsevents

import "errors"

type UXFlowCompletedV1 struct {
	FlowID       string `json:"flow_id"`
	ParentUserID string `json:"parent_user_id"`
	ChildProfile string `json:"child_profile_id"`
	Step         string `json:"step"`
	Status       string `json:"status"`
	CompletedAt  string `json:"completed_at"`
}

func (e UXFlowCompletedV1) Validate() error {
	if e.FlowID == "" || e.ParentUserID == "" || e.Step == "" || e.Status == "" || e.CompletedAt == "" {
		return errors.New("ux.flow.completed.v1 has missing required fields")
	}
	return nil
}
