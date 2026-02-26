package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetParentControls(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childProfileID := r.PathValue("child_profile_id")
		if childProfileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_child_profile", "child_profile_id is required")
			return
		}
		parentID, err := parentUserIDFromRequest(r, r.URL.Query().Get("parent_user_id"))
		if err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "parent_mismatch", err.Error())
			return
		}
		if parentID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_parent", "parent_user_id is required")
			return
		}
		if err := ensureChildOwnership(r, parentID, childProfileID); err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
			return
		}
		controls, err := store.GetControls(childProfileID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "controls_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ParentControlsResponse{
			ChildProfileID: childProfileID,
			Controls:       controls,
		})
	}
}
