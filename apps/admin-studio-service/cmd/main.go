package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/admin-studio-service/internal"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
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
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	mux := internal.NewMuxWithBus(repo, bus)
	addr := ":8090"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("admin-studio-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
