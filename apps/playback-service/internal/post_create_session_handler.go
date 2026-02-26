package internal

import (
	"errors"
	"net/http"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

func PostCreateSession(service *Service, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.CreatePlaybackSessionRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if req.SafetyMode == "" {
			req.SafetyMode = contractsapi.SafetyModeStrict
		}
		if req.EntitlementStatus == "" {
			req.EntitlementStatus = "active"
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		result, err := service.CreateSession(r.Context(), req)
		if err != nil {
			if errors.Is(err, ErrEntitlementRequired) {
				httpx.WriteAPIError(w, http.StatusForbidden, "entitlement_required", err.Error())
				return
			}
			if errors.Is(err, ErrParentGateRequired) {
				httpx.WriteAPIError(w, http.StatusForbidden, "parent_gate_required", err.Error())
				return
			}
			if errors.Is(err, ErrChildProfileForbidden) {
				httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
				return
			}
			if errors.Is(err, ErrSessionCapReached) {
				httpx.WriteAPIError(w, http.StatusForbidden, "session_cap_reached", err.Error())
				return
			}
			if errors.Is(err, ErrSessionLimitTooHigh) {
				httpx.WriteAPIError(w, http.StatusBadRequest, "session_limit_too_high", err.Error())
				return
			}
			if errors.Is(err, ErrAutoplayBlocked) {
				httpx.WriteAPIError(w, http.StatusForbidden, "autoplay_blocked", err.Error())
				return
			}
			httpx.WriteAPIError(w, http.StatusInternalServerError, "playback_error", err.Error())
			return
		}
		if err := publishEvent(r.Context(), bus, "playback.session.started.v1", uuid.NewString(), contractsevents.PlaybackSessionStartedV1{
			PlaybackSessionID: result.PlaybackSessionID,
			ChildProfileID:    req.ChildProfileID,
			EpisodeID:         req.EpisodeID,
			StartedAt:         time.Now().UTC().Format(time.RFC3339),
			SafetyMode:        req.SafetyMode,
		}); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "event_publish_failed", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, result)
	}
}
