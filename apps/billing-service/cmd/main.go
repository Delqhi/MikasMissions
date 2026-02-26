package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/billing-service/internal"
)

func main() {
	repository, err := internal.OpenRepositoryFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if closer, ok := repository.(io.Closer); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("close repository: %v", err)
			}
		}()
	}
	mux := internal.NewMux(repository)
	addr := ":8089"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("billing-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
