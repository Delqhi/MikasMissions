package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PutAdminWorkflow(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID := r.PathValue("workflow_id")
		if workflowID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "workflow_id is required")
			return
		}
		var req contractsapi.UpdateAdminWorkflowRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		actor := "admin-system"
		if principal, ok := authz.PrincipalFrom(r.Context()); ok {
			actor = actorIDFromPrincipal(principal)
		}
		updated, found, err := repo.UpdateWorkflow(WorkflowTemplate{
			ID:                 workflowID,
			Name:               req.Name,
			Description:        req.Description,
			ContentSuitability: req.ContentSuitability,
			AgeBand:            req.AgeBand,
			Steps:              req.Steps,
			ModelProfileID:     req.ModelProfileID,
			SafetyProfile:      req.SafetyProfile,
		}, actor)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !found {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "workflow not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, mapWorkflowToContract(updated))
	}
}
