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

func PostAdminRunRetry(repo Repository, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID := r.PathValue("run_id")
		if runID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "run_id is required")
			return
		}
		run, found, err := repo.FindRun(runID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !found {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "run not found")
			return
		}
		workflow, found, err := repo.FindWorkflow(run.WorkflowID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !found {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "workflow not found")
			return
		}
		if _, err := repo.SetRunStatus(runID, "requested", ""); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		actor := "admin-system"
		if principal, ok := authz.PrincipalFrom(r.Context()); ok {
			actor = actorIDFromPrincipal(principal)
		}
		_ = repo.AppendRunLog(WorkflowRunLog{
			RunID:     runID,
			Step:      "run",
			Status:    "retry_requested",
			Message:   "workflow run retry requested",
			EventTime: time.Now().UTC().Format(time.RFC3339),
		})
		err = publishJSONEvent(r.Context(), bus, "video.run.requested.v1", contractsevents.VideoRunRequestedV1{
			RunID:              run.ID,
			WorkflowID:         run.WorkflowID,
			ModelProfileID:     workflow.ModelProfileID,
			InputPayload:       run.InputPayload,
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
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminWorkflowRunResponse{RunID: runID, Status: "requested"})
	}
}
