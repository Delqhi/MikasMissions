package internal

import "net/http"

func NewMux(upstreams Upstreams) http.Handler {
	mux := http.NewServeMux()
	identity := proxyTo(upstreams.Identity)
	profile := proxyTo(upstreams.Profile)
	catalog := proxyTo(upstreams.Catalog)
	recommendation := proxyTo(upstreams.Recommendation)
	playback := proxyTo(upstreams.Playback)
	progress := proxyTo(upstreams.Progress)
	creator := proxyTo(upstreams.Creator)
	billing := proxyTo(upstreams.Billing)
	adminStudio := proxyTo(upstreams.AdminStudio)

	mux.Handle("POST /v1/parents/signup", identity)
	mux.Handle("POST /v1/parents/login", identity)
	mux.Handle("POST /v1/admin/login", identity)
	mux.Handle("POST /v1/parents/consent/verify", identity)
	mux.Handle("GET /v1/parents/dashboard", identity)
	mux.Handle("GET /v1/parents/controls/{child_profile_id}", identity)
	mux.Handle("PUT /v1/parents/controls/{child_profile_id}", identity)
	mux.Handle("POST /v1/parents/gates/challenge", identity)
	mux.Handle("POST /v1/parents/gates/verify", identity)

	mux.Handle("POST /v1/children/profiles", profile)
	mux.Handle("GET /v1/children/profiles", profile)

	mux.Handle("GET /v1/home/rails", recommendation)
	mux.Handle("GET /v1/kids/home", recommendation)
	mux.Handle("GET /v1/kids/progress/{child_profile_id}", progress)

	mux.Handle("GET /v1/catalog/episodes/{id}", catalog)

	mux.Handle("POST /v1/playback/sessions", playback)

	mux.Handle("POST /v1/progress/watch-events", progress)

	mux.Handle("GET /v1/billing/entitlements", billing)

	mux.Handle("POST /v1/creator/assets/upload", creator)
	mux.Handle("GET /v1/admin/workflows", adminStudio)
	mux.Handle("POST /v1/admin/workflows", adminStudio)
	mux.Handle("PUT /v1/admin/workflows/{workflow_id}", adminStudio)
	mux.Handle("DELETE /v1/admin/workflows/{workflow_id}", adminStudio)
	mux.Handle("POST /v1/admin/workflows/{workflow_id}/runs", adminStudio)
	mux.Handle("GET /v1/admin/runs/{run_id}", adminStudio)
	mux.Handle("GET /v1/admin/runs/{run_id}/logs", adminStudio)
	mux.Handle("POST /v1/admin/runs/{run_id}/retry", adminStudio)
	mux.Handle("POST /v1/admin/runs/{run_id}/cancel", adminStudio)
	mux.Handle("GET /v1/admin/model-profiles/{id}", adminStudio)
	mux.Handle("PUT /v1/admin/model-profiles/{id}", adminStudio)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return newGatewayHandler(mux, newGatewayAuthFromEnv())
}
