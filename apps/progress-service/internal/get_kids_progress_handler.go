package internal

import (
	"errors"
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetKidsProgress(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childProfileID := r.PathValue("child_profile_id")
		if childProfileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_child_profile", "child_profile_id is required")
			return
		}
		result, err := service.GetKidsProgress(r.Context(), childProfileID)
		if err != nil {
			if errors.Is(err, ErrChildProfileForbidden) {
				httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
				return
			}
			httpx.WriteAPIError(w, http.StatusInternalServerError, "progress_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, result)
	}
}
