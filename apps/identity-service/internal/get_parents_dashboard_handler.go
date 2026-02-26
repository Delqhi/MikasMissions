package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetParentsDashboard(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentID, err := parentUserIDFromRequest(r, r.URL.Query().Get("parent_user_id"))
		if err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "parent_mismatch", err.Error())
			return
		}
		if parentID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_parent", "parent_user_id is required")
			return
		}
		exists, err := store.ParentExists(parentID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "dashboard_error", err.Error())
			return
		}
		if !exists {
			httpx.WriteAPIError(w, http.StatusNotFound, "parent_not_found", "parent_user_id not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ParentsDashboardResponse{
			ParentUserID:   parentID,
			ChildProfiles:  0,
			WeeklyWatchMS:  0,
			WeeklyLearning: 0,
		})
	}
}
