package contractsevents

import "errors"

type VideoRunStepCompletedV1 struct {
	RunID       string `json:"run_id"`
	Step        string `json:"step"`
	Status      string `json:"status"`
	Details     string `json:"details"`
	CompletedAt string `json:"completed_at"`
}

func (e VideoRunStepCompletedV1) Validate() error {
	if e.RunID == "" || e.Step == "" || e.Status == "" || e.CompletedAt == "" {
		return errors.New("video.run.step.completed.v1 has missing required fields")
	}
	return nil
}
