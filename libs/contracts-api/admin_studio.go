package contractsapi

import (
	"encoding/json"
	"strings"
)

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AdminUserID string `json:"admin_user_id"`
	Role        string `json:"role"`
}

type AdminWorkflow struct {
	WorkflowID         string   `json:"workflow_id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	ContentSuitability string   `json:"content_suitability"`
	AgeBand            string   `json:"age_band"`
	Steps              []string `json:"steps"`
	ModelProfileID     string   `json:"model_profile_id"`
	SafetyProfile      string   `json:"safety_profile"`
	Version            int      `json:"version"`
}

type AdminWorkflowListResponse struct {
	Workflows []AdminWorkflow `json:"workflows"`
}

type CreateAdminWorkflowRequest struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	ContentSuitability string   `json:"content_suitability"`
	AgeBand            string   `json:"age_band"`
	Steps              []string `json:"steps"`
	ModelProfileID     string   `json:"model_profile_id"`
	SafetyProfile      string   `json:"safety_profile"`
}

type UpdateAdminWorkflowRequest struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	ContentSuitability string   `json:"content_suitability"`
	AgeBand            string   `json:"age_band"`
	Steps              []string `json:"steps"`
	ModelProfileID     string   `json:"model_profile_id"`
	SafetyProfile      string   `json:"safety_profile"`
}

type AdminWorkflowRunRequest struct {
	InputPayload json.RawMessage `json:"input_payload"`
	Priority     string          `json:"priority"`
	AutoPublish  bool            `json:"auto_publish"`
}

type AdminWorkflowRunResponse struct {
	RunID   string `json:"run_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type AdminRunResponse struct {
	RunID        string          `json:"run_id"`
	WorkflowID   string          `json:"workflow_id"`
	Status       string          `json:"status"`
	Priority     string          `json:"priority"`
	AutoPublish  bool            `json:"auto_publish"`
	InputPayload json.RawMessage `json:"input_payload"`
	LastError    string          `json:"last_error"`
}

type AdminRunLogEntry struct {
	RunID     string `json:"run_id"`
	Step      string `json:"step"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	EventTime string `json:"event_time"`
}

type AdminRunLogsResponse struct {
	RunID string             `json:"run_id"`
	Logs  []AdminRunLogEntry `json:"logs"`
}

type AdminModelProfile struct {
	ModelProfileID string `json:"model_profile_id"`
	Provider       string `json:"provider"`
	BaseURL        string `json:"base_url"`
	ModelID        string `json:"model_id"`
	TimeoutMS      int    `json:"timeout_ms"`
	MaxRetries     int    `json:"max_retries"`
	SafetyPreset   string `json:"safety_preset"`
}

func (r AdminLoginRequest) Validate() *APIError {
	if !strings.Contains(r.Email, "@") {
		return &APIError{Code: "invalid_email", Message: "email is invalid"}
	}
	if len(r.Password) < 10 {
		return &APIError{Code: "weak_password", Message: "password must have at least 10 chars"}
	}
	return nil
}

func (r CreateAdminWorkflowRequest) Validate() *APIError {
	if r.Name == "" || r.ModelProfileID == "" || len(r.Steps) == 0 {
		return &APIError{Code: "workflow_invalid", Message: "name, model_profile_id and steps are required"}
	}
	if !IsValidAgeBand(r.AgeBand) {
		return &APIError{Code: "workflow_invalid", Message: "age_band must be one of: 3-5, 6-11, 12-16"}
	}
	if r.ContentSuitability == "" {
		return &APIError{Code: "workflow_invalid", Message: "content_suitability is required"}
	}
	if r.SafetyProfile == "" {
		return &APIError{Code: "workflow_invalid", Message: "safety_profile is required"}
	}
	return nil
}

func (r UpdateAdminWorkflowRequest) Validate() *APIError {
	return CreateAdminWorkflowRequest(r).Validate()
}

func (r AdminWorkflowRunRequest) Normalize() AdminWorkflowRunRequest {
	if r.Priority == "" {
		r.Priority = "normal"
	}
	if len(r.InputPayload) == 0 {
		r.InputPayload = json.RawMessage(`{}`)
	}
	return r
}

func (p AdminModelProfile) Validate() *APIError {
	if p.ModelProfileID == "" || p.Provider == "" || p.BaseURL == "" || p.ModelID == "" {
		return &APIError{Code: "workflow_invalid", Message: "model_profile_id, provider, base_url and model_id are required"}
	}
	if p.Provider != "nvidia_nim" {
		return &APIError{Code: "workflow_invalid", Message: "provider must be nvidia_nim"}
	}
	if p.TimeoutMS < 500 || p.TimeoutMS > 120000 {
		return &APIError{Code: "workflow_invalid", Message: "timeout_ms must be between 500 and 120000"}
	}
	if p.MaxRetries < 0 || p.MaxRetries > 10 {
		return &APIError{Code: "workflow_invalid", Message: "max_retries must be between 0 and 10"}
	}
	if p.SafetyPreset == "" {
		return &APIError{Code: "workflow_invalid", Message: "safety_preset is required"}
	}
	return nil
}
