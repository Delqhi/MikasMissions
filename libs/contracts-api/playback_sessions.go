package contractsapi

import "time"

type CreatePlaybackSessionRequest struct {
	ChildProfileID      string `json:"child_profile_id"`
	EpisodeID           string `json:"episode_id"`
	DeviceType          string `json:"device_type"`
	SafetyMode          string `json:"safety_mode"`
	ParentGateToken     string `json:"parent_gate_token"`
	SessionLimitMinutes int    `json:"session_limit_minutes"`
	SessionMinutesUsed  int    `json:"session_minutes_used"`
	EntitlementStatus   string `json:"entitlement_status"`
	AutoplayRequested   bool   `json:"autoplay_requested"`
}

type CreatePlaybackSessionResponse struct {
	PlaybackSessionID string    `json:"playback_session_id"`
	Token             string    `json:"token"`
	StreamURL         string    `json:"stream_url"`
	ExpiresAt         time.Time `json:"expires_at"`
	SafetyApplied     bool      `json:"safety_applied"`
	SafetyReason      string    `json:"safety_reason"`
	SessionCapped     bool      `json:"session_capped"`
	SessionMaxMinutes int       `json:"session_max_minutes"`
}

func (r CreatePlaybackSessionRequest) Validate() *APIError {
	if r.ChildProfileID == "" || r.EpisodeID == "" || r.DeviceType == "" {
		return &APIError{Code: "missing_fields", Message: "child_profile_id, episode_id and device_type are required"}
	}
	if r.SafetyMode == "" {
		r.SafetyMode = SafetyModeStrict
	}
	if !ValidateSafetyMode(r.SafetyMode) {
		return &APIError{Code: "invalid_safety_mode", Message: "safety_mode must be strict or balanced"}
	}
	if r.SessionLimitMinutes != 0 && (r.SessionLimitMinutes < 5 || r.SessionLimitMinutes > 180) {
		return &APIError{Code: "invalid_session_limit", Message: "session_limit_minutes must be between 5 and 180"}
	}
	if r.EntitlementStatus == "" {
		r.EntitlementStatus = "active"
	}
	if r.EntitlementStatus != "active" && r.EntitlementStatus != "inactive" {
		return &APIError{Code: "invalid_entitlement_status", Message: "entitlement_status must be active or inactive"}
	}
	if r.SessionMinutesUsed < 0 {
		return &APIError{Code: "invalid_session_usage", Message: "session_minutes_used must be >= 0"}
	}
	return nil
}
