package contractsapi

import "strings"

type ParentSignupRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	Marketing  bool   `json:"marketing"`
	AcceptedTo bool   `json:"accepted_terms"`
}

type ParentSignupResponse struct {
	ParentUserID string `json:"parent_user_id"`
	Status       string `json:"status"`
}

func (r ParentSignupRequest) Validate() *APIError {
	if !r.AcceptedTo {
		return &APIError{Code: "terms_required", Message: "accepted_terms must be true"}
	}
	if !strings.Contains(r.Email, "@") {
		return &APIError{Code: "invalid_email", Message: "email is invalid"}
	}
	if len(r.Password) < 10 {
		return &APIError{Code: "weak_password", Message: "password must have at least 10 chars"}
	}
	if r.Country == "" || r.Language == "" {
		return &APIError{Code: "missing_profile", Message: "country and language are required"}
	}
	return nil
}
