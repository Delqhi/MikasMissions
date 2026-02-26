package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func GetEpisode(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		episodeID := r.PathValue("id")
		if episodeID == "" {
			httpx.WriteAPIError(w, http.StatusBadRequest, "missing_id", "episode id is required")
			return
		}
		episode, ok, err := repo.FindEpisode(episodeID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "episode_lookup_error", err.Error())
			return
		}
		if !ok {
			httpx.WriteAPIError(w, http.StatusNotFound, "episode_not_found", "episode id not found")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, episode)
	}
}
