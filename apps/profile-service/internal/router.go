package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
)

func NewMux(store Repository) *http.ServeMux {
	mux := http.NewServeMux()
	authorizer := authz.NewHTTPAuthorizerFromEnv()
	mux.HandleFunc("POST /v1/children/profiles", authorizer.Wrap([]string{"parent", "service"}, PostCreateProfile(store)))
	mux.HandleFunc("GET /v1/children/profiles", authorizer.Wrap([]string{"parent", "service"}, GetProfiles(store)))
	mux.HandleFunc("GET /internal/profiles/{child_profile_id}", GetInternalProfile(store))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
