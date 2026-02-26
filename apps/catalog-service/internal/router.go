package internal

import "net/http"

func NewMux(repo Repository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/catalog/episodes/{id}", GetEpisode(repo))
	mux.HandleFunc("GET /internal/catalog/episodes", GetInternalEpisodes(repo))
	mux.HandleFunc("POST /internal/catalog/episodes", PostInternalUpsertEpisode(repo))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
