package internal

import (
	"errors"
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetEntitlement(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentID := r.URL.Query().Get("parent_user_id")
		childID := r.URL.Query().Get("child_profile_id")
		if parentID == "" && childID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_subject", "parent_user_id or child_profile_id is required")
			return
		}

		var result EntitlementResponse
		var err error
		if parentID != "" {
			result, err = repo.GetEntitlementByParent(r.Context(), parentID)
		} else {
			result, err = repo.GetEntitlementByChild(r.Context(), childID)
		}
		if errors.Is(err, ErrEntitlementNotFound) {
			httpx.WriteAPIError(w, http.StatusNotFound, "entitlement_not_found", err.Error())
			return
		}
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "entitlement_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, result)
	}
}
