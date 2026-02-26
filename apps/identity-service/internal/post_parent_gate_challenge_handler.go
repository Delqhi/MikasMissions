package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostParentGateChallenge(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.ParentGateChallengeRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		parentUserID, err := parentUserIDFromRequest(r, req.ParentUserID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "parent_mismatch", err.Error())
			return
		}
		req.ParentUserID = parentUserID
		exists, err := store.ParentExists(req.ParentUserID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_challenge_error", err.Error())
			return
		}
		if !exists {
			httpx.WriteAPIError(w, http.StatusNotFound, "parent_not_found", "parent_user_id not found")
			return
		}
		if err := ensureChildOwnership(r, req.ParentUserID, req.ChildProfileID); err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
			return
		}
		challenge, err := store.CreateGateChallenge(req.ParentUserID, req.ChildProfileID, req.Method)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_challenge_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, contractsapi.ParentGateChallengeResponse{
			ChallengeID: challenge.ChallengeID,
			Method:      challenge.Method,
			ExpiresAt:   challenge.ExpiresAt,
		})
	}
}
