package contractsapi

type ParentConsentVerifyRequest struct {
	ParentUserID string `json:"parent_user_id"`
	Method       string `json:"method"`
	Challenge    string `json:"challenge"`
}

type ParentConsentVerifyResponse struct {
	ConsentID string `json:"consent_id"`
	Verified  bool   `json:"verified"`
}

func (r ParentConsentVerifyRequest) Validate() *APIError {
	if r.ParentUserID == "" {
		return &APIError{Code: "parent_required", Message: "parent_user_id is required"}
	}
	if r.Method == "" || r.Challenge == "" {
		return &APIError{Code: "verification_required", Message: "method and challenge are required"}
	}
	return nil
}
