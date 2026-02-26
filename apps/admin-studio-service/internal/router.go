package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func NewMux(repo Repository) *http.ServeMux {
	return NewMuxWithBus(repo, nil)
}

func NewMuxWithBus(repo Repository, bus queue.Bus) *http.ServeMux {
	authorizer := authz.NewHTTPAuthorizerFromEnv()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/admin/workflows", authorizer.Wrap([]string{"admin", "service"}, GetAdminWorkflows(repo)))
	mux.HandleFunc("POST /v1/admin/workflows", authorizer.Wrap([]string{"admin", "service"}, PostAdminWorkflow(repo, bus)))
	mux.HandleFunc("PUT /v1/admin/workflows/{workflow_id}", authorizer.Wrap([]string{"admin", "service"}, PutAdminWorkflow(repo)))
	mux.HandleFunc("DELETE /v1/admin/workflows/{workflow_id}", authorizer.Wrap([]string{"admin", "service"}, DeleteAdminWorkflow(repo)))
	mux.HandleFunc("POST /v1/admin/workflows/{workflow_id}/runs", authorizer.Wrap([]string{"admin", "service"}, PostAdminWorkflowRun(repo, bus)))
	mux.HandleFunc("GET /v1/admin/runs/{run_id}", authorizer.Wrap([]string{"admin", "service"}, GetAdminRun(repo)))
	mux.HandleFunc("GET /v1/admin/runs/{run_id}/logs", authorizer.Wrap([]string{"admin", "service"}, GetAdminRunLogs(repo)))
	mux.HandleFunc("POST /v1/admin/runs/{run_id}/retry", authorizer.Wrap([]string{"admin", "service"}, PostAdminRunRetry(repo, bus)))
	mux.HandleFunc("POST /v1/admin/runs/{run_id}/cancel", authorizer.Wrap([]string{"admin", "service"}, PostAdminRunCancel(repo)))
	mux.HandleFunc("GET /v1/admin/model-profiles/{id}", authorizer.Wrap([]string{"admin", "service"}, GetAdminModelProfile(repo)))
	mux.HandleFunc("PUT /v1/admin/model-profiles/{id}", authorizer.Wrap([]string{"admin", "service"}, PutAdminModelProfile(repo)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
