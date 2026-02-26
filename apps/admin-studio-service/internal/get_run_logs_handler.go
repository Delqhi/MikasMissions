package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetAdminRunLogs(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runID := r.PathValue("run_id")
		if runID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "run_id is required")
			return
		}
		logs, err := repo.ListRunLogs(runID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		mapped := make([]contractsapi.AdminRunLogEntry, 0, len(logs))
		for _, item := range logs {
			mapped = append(mapped, contractsapi.AdminRunLogEntry{
				RunID:     item.RunID,
				Step:      item.Step,
				Status:    item.Status,
				Message:   item.Message,
				EventTime: item.EventTime,
			})
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminRunLogsResponse{RunID: runID, Logs: mapped})
	}
}
