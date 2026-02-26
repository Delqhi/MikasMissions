package main

import (
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/recommendation-service/internal"
)

func main() {
	service := internal.NewService()
	mux := internal.NewMux(service)
	addr := ":8084"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("recommendation-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
