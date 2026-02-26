package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

func PutParentControls(store Repository, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childProfileID := r.PathValue("child_profile_id")
		if childProfileID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_child_profile", "child_profile_id is required")
			return
		}
		var controls contractsapi.ParentalControls
		if err := httpx.DecodeJSON(r, &controls); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if controls.SafetyMode == "" {
			controls.SafetyMode = contractsapi.SafetyModeStrict
		}
		if controls.SessionLimitMinutes == 0 {
			controls.SessionLimitMinutes = 30
		}
		if controls.BedtimeWindow == "" {
			controls.BedtimeWindow = "20:00-07:00"
		}
		if apiErr := controls.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		parentUserID, err := parentUserIDFromRequest(r, r.URL.Query().Get("parent_user_id"))
		if err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "parent_mismatch", err.Error())
			return
		}
		if parentUserID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_parent", "parent_user_id is required")
			return
		}
		if err := ensureChildOwnership(r, parentUserID, childProfileID); err != nil {
			httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
			return
		}
		if err := store.SetControls(childProfileID, controls); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "controls_error", err.Error())
			return
		}
		if err := publishEvent(r.Context(), bus, "parent.controls.updated.v1", uuid.NewString(), contractsevents.ParentControlsUpdatedV1{
			ParentUserID:     parentUserID,
			ChildProfileID:   childProfileID,
			SafetyMode:       controls.SafetyMode,
			SessionLimitMins: controls.SessionLimitMinutes,
			ExternalLinks:    controls.ExternalLinks,
			AuditEventID:     uuid.NewString(),
		}); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "event_publish_failed", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ParentControlsResponse{
			ChildProfileID: childProfileID,
			Controls:       controls,
		})
	}
}
