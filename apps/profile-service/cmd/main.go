package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/profile-service/internal"
)

func main() {
	store, err := internal.OpenRepositoryFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if closer, ok := store.(io.Closer); ok {
		defer func() {
			_ = closer.Close()
		}()
	}
	mux := internal.NewMux(store)
	addr := ":8082"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("profile-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
