package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func DeleteAdminWorkflow(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID := r.PathValue("workflow_id")
		if workflowID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "workflow_id is required")
			return
		}
		actor := "admin-system"
		if principal, ok := authz.PrincipalFrom(r.Context()); ok {
			actor = actorIDFromPrincipal(principal)
		}
		deleted, err := repo.DeleteWorkflow(workflowID, actor)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !deleted {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "workflow not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
