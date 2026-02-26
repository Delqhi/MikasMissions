package contractsevents

import "errors"

type VideoRunFailedV1 struct {
	RunID        string `json:"run_id"`
	Step         string `json:"step"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	FailedAt     string `json:"failed_at"`
}

func (e VideoRunFailedV1) Validate() error {
	if e.RunID == "" || e.Step == "" || e.ErrorCode == "" || e.ErrorMessage == "" || e.FailedAt == "" {
		return errors.New("video.run.failed.v1 has missing required fields")
	}
	return nil
}
