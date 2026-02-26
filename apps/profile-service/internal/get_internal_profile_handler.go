package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

type InternalProfileResponse struct {
	ChildProfileID string `json:"child_profile_id"`
	ParentUserID   string `json:"parent_user_id"`
	AgeBand        string `json:"age_band"`
}

func GetInternalProfile(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childProfileID := r.PathValue("child_profile_id")
		if childProfileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_child_profile", "child_profile_id is required")
			return
		}
		profile, ok, err := store.FindProfile(childProfileID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "profile_error", err.Error())
			return
		}
		if !ok {
			httpx.WriteAPIError(w, http.StatusNotFound, "profile_not_found", "child_profile_id not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, InternalProfileResponse{
			ChildProfileID: profile.ID,
			ParentUserID:   profile.ParentUser,
			AgeBand:        profile.AgeBand,
		})
	}
}
