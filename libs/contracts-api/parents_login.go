package contractsapi

import "strings"

type ParentLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ParentLoginResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ParentUserID string `json:"parent_user_id"`
	Role         string `json:"role"`
}

func (r ParentLoginRequest) Validate() *APIError {
	if !strings.Contains(r.Email, "@") {
		return &APIError{Code: "invalid_email", Message: "email is invalid"}
	}
	if len(r.Password) < 10 {
		return &APIError{Code: "weak_password", Message: "password must have at least 10 chars"}
	}
	return nil
}
