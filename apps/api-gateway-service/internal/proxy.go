package internal

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyTo(target *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.ErrorHandler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if original != nil {
			original(w, r, err)
			return
		}
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}
	return proxy
}
