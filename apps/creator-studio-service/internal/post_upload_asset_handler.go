package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostUploadAsset(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UploadRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if req.SourceURL == "" || req.UploaderID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_fields", "source_url and uploader_id are required")
			return
		}
		result, err := service.UploadAsset(r.Context(), req)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "upload_failed", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusAccepted, result)
	}
}
