package contractsapi

const (
	AgeBandEarly = "3-5"
	AgeBandCore  = "6-11"
	AgeBandTeen  = "12-16"
)

type CreateChildProfileRequest struct {
	ParentUserID string `json:"parent_user_id"`
	DisplayName  string `json:"display_name"`
	AgeBand      string `json:"age_band"`
	Avatar       string `json:"avatar"`
}

type CreateChildProfileResponse struct {
	ChildProfileID string `json:"child_profile_id"`
	SafeDefaults   bool   `json:"safe_defaults"`
}

type ChildProfileSummary struct {
	ChildProfileID string `json:"child_profile_id"`
	ParentUserID   string `json:"parent_user_id"`
	DisplayName    string `json:"display_name"`
	AgeBand        string `json:"age_band"`
	Avatar         string `json:"avatar"`
}

type ListChildProfilesResponse struct {
	Profiles []ChildProfileSummary `json:"profiles"`
}

func (r CreateChildProfileRequest) Validate() *APIError {
	if r.ParentUserID == "" || r.DisplayName == "" {
		return &APIError{Code: "missing_fields", Message: "parent_user_id and display_name are required"}
	}
	if !IsValidAgeBand(r.AgeBand) {
		return &APIError{Code: "invalid_age_band", Message: "age_band must be one of: 3-5, 6-11, 12-16"}
	}
	return nil
}

func IsValidAgeBand(ageBand string) bool {
	switch ageBand {
	case AgeBandEarly, AgeBandCore, AgeBandTeen:
		return true
	default:
		return false
	}
}
