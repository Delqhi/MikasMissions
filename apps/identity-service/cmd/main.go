package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/identity-service/internal"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
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
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	mux := internal.NewMuxWithBus(store, bus)
	addr := ":8081"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("identity-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
