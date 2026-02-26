package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetAdminModelProfile(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := r.PathValue("id")
		if profileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "model profile id is required")
			return
		}
		profile, found, err := repo.GetModelProfile(profileID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
			return
		}
		if !found {
			httpx.WriteAPIError(w, http.StatusNotFound, "workflow_missing", "model profile not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminModelProfile{
			ModelProfileID: profile.ID,
			Provider:       profile.Provider,
			BaseURL:        profile.BaseURL,
			ModelID:        profile.ModelID,
			TimeoutMS:      profile.TimeoutMS,
			MaxRetries:     profile.MaxRetries,
			SafetyPreset:   profile.SafetyPreset,
		})
	}
}
