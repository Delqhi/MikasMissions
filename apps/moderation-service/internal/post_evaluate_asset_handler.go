package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

type EvaluateAssetRequest struct {
	AssetID string `json:"asset_id"`
	AgeBand string `json:"age_band"`
}

type EvaluateAssetResponse struct {
	PolicyResult string `json:"policy_result"`
}

func PostEvaluateAsset(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EvaluateAssetRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if req.AssetID == "" || req.AgeBand == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_fields", "asset_id and age_band are required")
			return
		}
		result, err := service.EvaluateAsset(r.Context(), req.AssetID, req.AgeBand)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "moderation_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, EvaluateAssetResponse{PolicyResult: result})
	}
}
