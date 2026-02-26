package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PutAdminModelProfile(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := r.PathValue("id")
		if profileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "workflow_invalid", "model profile id is required")
			return
		}
		var req contractsapi.AdminModelProfile
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		req.ModelProfileID = profileID
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		actor := "admin-system"
		if principal, ok := authz.PrincipalFrom(r.Context()); ok {
			actor = actorIDFromPrincipal(principal)
		}
		profile, err := repo.PutModelProfile(ModelProfile{
			ID:           req.ModelProfileID,
			Provider:     req.Provider,
			BaseURL:      req.BaseURL,
			ModelID:      req.ModelID,
			TimeoutMS:    req.TimeoutMS,
			MaxRetries:   req.MaxRetries,
			SafetyPreset: req.SafetyPreset,
		}, actor)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "workflow_error", err.Error())
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
