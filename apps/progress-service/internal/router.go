package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func NewMux(service *Service) *http.ServeMux {
	return NewMuxWithBus(service, nil)
}

func NewMuxWithBus(service *Service, bus queue.Bus) *http.ServeMux {
	mux := http.NewServeMux()
	authorizer := authz.NewHTTPAuthorizerFromEnv()
	mux.HandleFunc("POST /v1/progress/watch-events", authorizer.Wrap([]string{"parent", "child", "service"}, PostWatchEvent(service, bus)))
	mux.HandleFunc("GET /v1/kids/progress/{child_profile_id}", authorizer.Wrap([]string{"parent", "child", "service"}, GetKidsProgress(service)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
