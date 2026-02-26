package contractsevents

import "errors"

type ChildSessionCappedV1 struct {
	ChildProfileID string `json:"child_profile_id"`
	SessionID      string `json:"session_id"`
	CappedAt       string `json:"capped_at"`
	Reason         string `json:"reason"`
	LimitMinutes   int    `json:"limit_minutes"`
}

func (e ChildSessionCappedV1) Validate() error {
	if e.ChildProfileID == "" || e.SessionID == "" || e.CappedAt == "" || e.Reason == "" {
		return errors.New("child.session.capped.v1 has missing required fields")
	}
	if e.LimitMinutes <= 0 {
		return errors.New("child.session.capped.v1 limit_minutes must be > 0")
	}
	return nil
}
