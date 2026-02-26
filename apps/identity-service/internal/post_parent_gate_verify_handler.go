package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"github.com/google/uuid"
)

func PostParentGateVerify(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.ParentGateVerifyRequest
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
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_verify_error", err.Error())
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
		consumed, err := store.ConsumeGateChallenge(req.ChallengeID, req.ParentUserID, req.ChildProfileID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_verify_error", err.Error())
			return
		}
		if !consumed {
			httpx.WriteAPIError(w, http.StatusForbidden, "invalid_or_replayed_challenge", "challenge is invalid, expired or already used")
			return
		}
		gateToken := uuid.NewString()
		if err := store.SaveGateToken(req.ChildProfileID, gateToken); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_verify_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ParentGateVerifyResponse{Verified: true, GateToken: gateToken})
	}
}
