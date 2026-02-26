package contractsevents

import "errors"

type ParentControlsUpdatedV1 struct {
	ParentUserID     string `json:"parent_user_id"`
	ChildProfileID   string `json:"child_profile_id"`
	SafetyMode       string `json:"safety_mode"`
	SessionLimitMins int    `json:"session_limit_minutes"`
	ExternalLinks    bool   `json:"external_links"`
	AuditEventID     string `json:"audit_event_id"`
}

func (e ParentControlsUpdatedV1) Validate() error {
	if e.ParentUserID == "" || e.ChildProfileID == "" || e.SafetyMode == "" || e.AuditEventID == "" {
		return errors.New("parent.controls.updated.v1 has missing required fields")
	}
	if e.SessionLimitMins <= 0 {
		return errors.New("parent.controls.updated.v1 requires positive session_limit_minutes")
	}
	return nil
}
