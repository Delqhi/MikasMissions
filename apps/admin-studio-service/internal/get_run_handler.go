package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetAdminRun(repo Repository) http.HandlerFunc {
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
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminRunResponse{
			RunID:        run.ID,
			WorkflowID:   run.WorkflowID,
			Status:       run.Status,
			Priority:     run.Priority,
			AutoPublish:  run.AutoPublish,
			InputPayload: run.InputPayload,
			LastError:    run.LastError,
		})
	}
}
