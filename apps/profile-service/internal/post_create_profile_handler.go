package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostCreateProfile(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.CreateChildProfileRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if principal, ok := authz.PrincipalFrom(r.Context()); ok && principal.Role == "parent" {
			if principal.ParentUserID == "" {
				httpx.WriteAPIError(w, http.StatusUnauthorized, "missing_identity", "parent identity is required")
				return
			}
			req.ParentUserID = principal.ParentUserID
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		profile, err := store.CreateProfile(req.ParentUserID, req.DisplayName, req.AgeBand, req.Avatar)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "profile_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, contractsapi.CreateChildProfileResponse{
			ChildProfileID: profile.ID,
			SafeDefaults:   true,
		})
	}
}
