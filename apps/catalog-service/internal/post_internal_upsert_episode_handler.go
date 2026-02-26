package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostInternalUpsertEpisode(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EpisodeUpsertRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if req.EpisodeID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_episode_id", "episode_id is required")
			return
		}
		if req.AgeBand != "" && !contractsapi.IsValidAgeBand(req.AgeBand) {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_age_band", "age_band must be 3-5, 6-11, or 12-16")
			return
		}
		if err := repo.UpsertEpisode(req); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "upsert_episode_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
