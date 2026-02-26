package internal

type EntitlementResponse struct {
	ParentUserID string `json:"parent_user_id"`
	Plan         string `json:"plan"`
	Active       bool   `json:"active"`
}
