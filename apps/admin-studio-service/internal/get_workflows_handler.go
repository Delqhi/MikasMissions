package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetAdminWorkflows(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		workflows, err := repo.ListWorkflows()
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		mapped := make([]contractsapi.AdminWorkflow, 0, len(workflows))
		for _, workflow := range workflows {
			mapped = append(mapped, mapWorkflowToContract(workflow))
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminWorkflowListResponse{Workflows: mapped})
	}
}

func mapWorkflowToContract(workflow WorkflowTemplate) contractsapi.AdminWorkflow {
	return contractsapi.AdminWorkflow{
		WorkflowID:         workflow.ID,
		Name:               workflow.Name,
		Description:        workflow.Description,
		ContentSuitability: workflow.ContentSuitability,
		AgeBand:            workflow.AgeBand,
		Steps:              workflow.Steps,
		ModelProfileID:     workflow.ModelProfileID,
		SafetyProfile:      workflow.SafetyProfile,
		Version:            workflow.Version,
	}
}
