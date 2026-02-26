package internal

import "net/http"

func NewMux(repo Repository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/billing/entitlements", GetEntitlement(repo))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
