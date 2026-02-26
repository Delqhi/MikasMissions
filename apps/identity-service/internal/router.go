package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func NewMux(store Repository) *http.ServeMux {
	return NewMuxWithBus(store, nil)
}

func NewMuxWithBus(store Repository, bus queue.Bus) *http.ServeMux {
	mux := http.NewServeMux()
	authorizer := authz.NewHTTPAuthorizerFromEnv()
	mux.HandleFunc("POST /v1/parents/signup", authorizer.Wrap(nil, PostSignup(store)))
	mux.HandleFunc("POST /v1/parents/login", authorizer.Wrap(nil, PostLogin(store)))
	mux.HandleFunc("POST /v1/admin/login", authorizer.Wrap(nil, PostAdminLogin(store)))
	mux.HandleFunc("POST /v1/parents/consent/verify", authorizer.Wrap(nil, PostConsentVerify(store, bus)))
	mux.HandleFunc("GET /v1/parents/dashboard", authorizer.Wrap([]string{"parent", "service"}, GetParentsDashboard(store)))
	mux.HandleFunc("GET /v1/parents/controls/{child_profile_id}", authorizer.Wrap([]string{"parent", "service"}, GetParentControls(store)))
	mux.HandleFunc("PUT /v1/parents/controls/{child_profile_id}", authorizer.Wrap([]string{"parent", "service"}, PutParentControls(store, bus)))
	mux.HandleFunc("POST /v1/parents/gates/challenge", authorizer.Wrap([]string{"parent", "service"}, PostParentGateChallenge(store)))
	mux.HandleFunc("POST /v1/parents/gates/verify", authorizer.Wrap([]string{"parent", "service"}, PostParentGateVerify(store)))
	mux.HandleFunc("POST /internal/gates/consume", PostInternalConsumeGateToken(store))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
