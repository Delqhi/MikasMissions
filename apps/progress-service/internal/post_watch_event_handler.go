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

func PostWatchEvent(service *Service, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.UpsertWatchEventRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		result, err := service.UpsertProgress(r.Context(), req)
		if err != nil {
			if errors.Is(err, ErrChildProfileForbidden) {
				httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
				return
			}
			httpx.WriteAPIError(w, http.StatusInternalServerError, "progress_error", err.Error())
			return
		}
		endedAt := req.EventTimeISO
		if _, err := time.Parse(time.RFC3339, endedAt); err != nil {
			endedAt = time.Now().UTC().Format(time.RFC3339)
		}
		progress, err := service.GetKidsProgress(r.Context(), req.ChildProfileID)
		if err != nil {
			if errors.Is(err, ErrChildProfileForbidden) {
				httpx.WriteAPIError(w, http.StatusForbidden, "child_profile_forbidden", err.Error())
				return
			}
			httpx.WriteAPIError(w, http.StatusInternalServerError, "progress_error", err.Error())
			return
		}
		if err := publishEvent(r.Context(), bus, "playback.session.ended.v1", uuid.NewString(), contractsevents.PlaybackSessionEndedV1{
			PlaybackSessionID: req.ChildProfileID + ":" + req.EpisodeID,
			ChildProfileID:    req.ChildProfileID,
			EpisodeID:         req.EpisodeID,
			EndedAt:           endedAt,
			WatchedMS:         req.WatchMS,
			Capped:            progress.SessionCapped,
		}); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "event_publish_failed", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusAccepted, result)
	}
}
