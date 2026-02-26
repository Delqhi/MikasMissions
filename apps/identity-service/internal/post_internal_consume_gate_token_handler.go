package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

type ConsumeGateTokenRequest struct {
	ChildProfileID string `json:"child_profile_id"`
	GateToken      string `json:"gate_token"`
}

type ConsumeGateTokenResponse struct {
	Consumed bool `json:"consumed"`
}

func PostInternalConsumeGateToken(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ConsumeGateTokenRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if req.ChildProfileID == "" || req.GateToken == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_fields", "child_profile_id and gate_token are required")
			return
		}
		consumed, err := store.ConsumeGateToken(req.ChildProfileID, req.GateToken)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "gate_token_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, ConsumeGateTokenResponse{Consumed: consumed})
	}
}
