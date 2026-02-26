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

func PostAdminWorkflowRun(repo Repository, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID := r.PathValue("workflow_id")
		if workflowID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "workflow_id is required")
			return
		}
		workflow, found, err := repo.FindWorkflow(workflowID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !found {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "workflow not found")
			return
		}
		var req contractsapi.AdminWorkflowRunRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		req = req.Normalize()
		actor := "admin-system"
		if principal, ok := authz.PrincipalFrom(r.Context()); ok {
			actor = actorIDFromPrincipal(principal)
		}
		run, err := repo.CreateRun(WorkflowRun{
			WorkflowID:   workflowID,
			Priority:     req.Priority,
			AutoPublish:  req.AutoPublish,
			InputPayload: req.InputPayload,
		}, actor)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		_ = repo.AppendRunLog(WorkflowRunLog{
			RunID:     run.ID,
			Step:      "run",
			Status:    "requested",
			Message:   "workflow run requested",
			EventTime: time.Now().UTC().Format(time.RFC3339),
		})
		err = publishJSONEvent(r.Context(), bus, "video.run.requested.v1", contractsevents.VideoRunRequestedV1{
			RunID:              run.ID,
			WorkflowID:         run.WorkflowID,
			ModelProfileID:     workflow.ModelProfileID,
			InputPayload:       req.InputPayload,
			AutoPublish:        run.AutoPublish,
			Priority:           run.Priority,
			ContentSuitability: workflow.ContentSuitability,
			AgeBand:            workflow.AgeBand,
			RequestedBy:        actor,
			RequestedAt:        time.Now().UTC().Format(time.RFC3339),
			TraceID:            run.ID,
		})
		if err != nil {
			httpx.WriteAPIError(w, http.StatusBadGateway, "nim_provider_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, contractsapi.AdminWorkflowRunResponse{
			RunID:  run.ID,
			Status: run.Status,
		})
	}
}
