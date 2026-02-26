package contractsevents

import "errors"

type ParentGateChallengeV1 struct {
	GateID        string `json:"gate_id"`
	ParentUserID  string `json:"parent_user_id"`
	ChildProfile  string `json:"child_profile_id"`
	Method        string `json:"method"`
	Verified      bool   `json:"verified"`
	ChallengeTime string `json:"challenge_time"`
}

func (e ParentGateChallengeV1) Validate() error {
	if e.GateID == "" || e.ParentUserID == "" || e.Method == "" || e.ChallengeTime == "" {
		return errors.New("parent.gate.challenge.v1 has missing required fields")
	}
	return nil
}
