package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/catalog-service/internal"
)

func main() {
	repo, err := internal.OpenRepositoryFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if closer, ok := repo.(io.Closer); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("close repository: %v", err)
			}
		}()
	}
	mux := internal.NewMux(repo)
	addr := ":8083"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("catalog-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
