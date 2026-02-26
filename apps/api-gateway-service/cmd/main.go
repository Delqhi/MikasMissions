package main

import (
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/api-gateway-service/internal"
)

func main() {
	upstreams, err := internal.LoadUpstreams()
	if err != nil {
		log.Fatal(err)
	}
	mux := internal.NewMux(upstreams)
	addr := ":8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("api-gateway-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
