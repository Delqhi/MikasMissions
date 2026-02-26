package internal

import (
	"net/http"
	"strconv"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

type listEpisodesResponse struct {
	Episodes []contractsapi.CatalogEpisodeResponse `json:"episodes"`
}

func GetInternalEpisodes(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ageBand := r.URL.Query().Get("age_band")
		limit := 20
		if raw := r.URL.Query().Get("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed <= 0 {
				httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_limit", "limit must be a positive integer")
				return
			}
			limit = parsed
		}
		episodes, err := repo.ListEpisodes(ageBand, limit)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "list_episodes_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, listEpisodesResponse{Episodes: episodes})
	}
}
