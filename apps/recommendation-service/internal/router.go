package internal

import (
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
)

func NewMux(service *Service) *http.ServeMux {
	mux := http.NewServeMux()
	authorizer := authz.NewHTTPAuthorizerFromEnv()
	mux.HandleFunc("GET /v1/home/rails", authorizer.Wrap([]string{"parent", "child", "service"}, GetHomeRails(service)))
	mux.HandleFunc("GET /v1/kids/home", authorizer.Wrap([]string{"parent", "child", "service"}, GetKidsHome(service)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
