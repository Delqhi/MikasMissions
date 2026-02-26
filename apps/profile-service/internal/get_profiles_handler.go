package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetProfiles(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentID := r.URL.Query().Get("parent_user_id")
		principal, ok := authz.PrincipalFrom(r.Context())
		if ok && principal.Role == "parent" {
			parentID = principal.ParentUserID
		}
		if parentID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_parent", "parent_user_id is required")
			return
		}
		profiles, err := store.ListProfilesByParent(parentID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "profiles_error", err.Error())
			return
		}
		out := make([]contractsapi.ChildProfileSummary, 0, len(profiles))
		for _, profile := range profiles {
			out = append(out, contractsapi.ChildProfileSummary{
				ChildProfileID: profile.ID,
				ParentUserID:   profile.ParentUser,
				DisplayName:    profile.DisplayName,
				AgeBand:        profile.AgeBand,
				Avatar:         profile.Avatar,
			})
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ListChildProfilesResponse{
			Profiles: out,
		})
	}
}
