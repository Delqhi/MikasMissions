package contractsevents

import "errors"

type ConsentVerifiedV1 struct {
	ConsentID    string `json:"consent_id"`
	ParentUserID string `json:"parent_user_id"`
	Method       string `json:"method"`
	VerifiedAt   string `json:"verified_at"`
}

func (e ConsentVerifiedV1) Validate() error {
	if e.ConsentID == "" || e.ParentUserID == "" || e.Method == "" || e.VerifiedAt == "" {
		return errors.New("consent.verified.v1 has missing required fields")
	}
	return nil
}
