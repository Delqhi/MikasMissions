package contractsapi

import "time"

const (
	SafetyModeStrict   = "strict"
	SafetyModeBalanced = "balanced"
)

type ParentalControls struct {
	Autoplay            bool   `json:"autoplay"`
	ChatEnabled         bool   `json:"chat_enabled"`
	ExternalLinks       bool   `json:"external_links"`
	SessionLimitMinutes int    `json:"session_limit_minutes"`
	BedtimeWindow       string `json:"bedtime_window"`
	SafetyMode          string `json:"safety_mode"`
}

type ParentControlsResponse struct {
	ChildProfileID string           `json:"child_profile_id"`
	Controls       ParentalControls `json:"controls"`
}

type ParentGateVerifyRequest struct {
	ParentUserID   string `json:"parent_user_id"`
	ChildProfileID string `json:"child_profile_id"`
	ChallengeID    string `json:"challenge_id"`
	Response       string `json:"response"`
}

type ParentGateVerifyResponse struct {
	Verified  bool   `json:"verified"`
	GateToken string `json:"gate_token"`
}

type ParentGateChallengeRequest struct {
	ParentUserID   string `json:"parent_user_id"`
	ChildProfileID string `json:"child_profile_id"`
	Method         string `json:"method"`
}

type ParentGateChallengeResponse struct {
	ChallengeID string    `json:"challenge_id"`
	Method      string    `json:"method"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func DefaultStrictControls() ParentalControls {
	return ParentalControls{
		Autoplay:            false,
		ChatEnabled:         false,
		ExternalLinks:       false,
		SessionLimitMinutes: 30,
		BedtimeWindow:       "20:00-07:00",
		SafetyMode:          SafetyModeStrict,
	}
}

func ValidateSafetyMode(mode string) bool {
	switch mode {
	case SafetyModeStrict, SafetyModeBalanced:
		return true
	default:
		return false
	}
}

func (c ParentalControls) Validate() *APIError {
	if !ValidateSafetyMode(c.SafetyMode) {
		return &APIError{Code: "invalid_safety_mode", Message: "safety_mode must be strict or balanced"}
	}
	if c.SessionLimitMinutes < 5 || c.SessionLimitMinutes > 180 {
		return &APIError{Code: "invalid_session_limit", Message: "session_limit_minutes must be between 5 and 180"}
	}
	if c.BedtimeWindow == "" {
		return &APIError{Code: "missing_bedtime", Message: "bedtime_window is required"}
	}
	return nil
}

func (r ParentGateVerifyRequest) Validate() *APIError {
	if r.ParentUserID == "" || r.ChildProfileID == "" {
		return &APIError{Code: "missing_fields", Message: "parent_user_id and child_profile_id are required"}
	}
	if r.ChallengeID == "" || r.Response == "" {
		return &APIError{Code: "missing_verification", Message: "challenge_id and response are required"}
	}
	return nil
}

func (r ParentGateChallengeRequest) Validate() *APIError {
	if r.ParentUserID == "" || r.ChildProfileID == "" || r.Method == "" {
		return &APIError{Code: "missing_fields", Message: "parent_user_id, child_profile_id and method are required"}
	}
	return nil
}
