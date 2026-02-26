package main

import (
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/moderation-service/internal"
)

func main() {
	service := internal.NewService()
	mux := internal.NewMux(service)
	addr := ":8088"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("moderation-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
