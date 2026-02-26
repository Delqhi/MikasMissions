package internal

import (
	"net/http"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func PostAdminWorkflow(repo Repository, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.CreateAdminWorkflowRequest
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
		created, err := repo.CreateWorkflow(WorkflowTemplate{
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
		_ = publishJSONEvent(r.Context(), bus, "video.workflow.created.v1", contractsevents.VideoWorkflowCreatedV1{
			WorkflowID: created.ID,
			Version:    created.Version,
			CreatedBy:  actor,
			CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		})
		httpx.WriteJSON(w, http.StatusCreated, mapWorkflowToContract(created))
	}
}
