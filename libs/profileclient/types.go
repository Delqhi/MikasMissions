package profileclient

import "errors"

var ErrProfileNotFound = errors.New("child profile not found")

type Profile struct {
	ChildProfileID string `json:"child_profile_id"`
	ParentUserID   string `json:"parent_user_id"`
	AgeBand        string `json:"age_band"`
}
