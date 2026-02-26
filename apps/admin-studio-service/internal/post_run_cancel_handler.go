package internal

import (
	"net/http"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostAdminRunCancel(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID := r.PathValue("run_id")
		if runID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "run_id is required")
			return
		}
		updated, err := repo.SetRunStatus(runID, "cancelled", "")
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !updated {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "run not found")
			return
		}
		_ = repo.AppendRunLog(WorkflowRunLog{
			RunID:     runID,
			Step:      "run",
			Status:    "cancelled",
			Message:   "workflow run cancelled by admin",
			EventTime: time.Now().UTC().Format(time.RFC3339),
		})
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminWorkflowRunResponse{RunID: runID, Status: "cancelled"})
	}
}
