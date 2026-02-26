package contractsevents

import "errors"

type SafetyFilterAppliedV1 struct {
	ChildProfileID string   `json:"child_profile_id"`
	EpisodeID      string   `json:"episode_id"`
	SafetyMode     string   `json:"safety_mode"`
	Filters        []string `json:"filters"`
	Reason         string   `json:"reason"`
}

func (e SafetyFilterAppliedV1) Validate() error {
	if e.ChildProfileID == "" || e.EpisodeID == "" || e.SafetyMode == "" || e.Reason == "" {
		return errors.New("safety.filter.applied.v1 has missing required fields")
	}
	if len(e.Filters) == 0 {
		return errors.New("safety.filter.applied.v1 requires at least one filter")
	}
	return nil
}
