package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/progress-service/internal"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
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

	profileReader, err := profileclient.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	service := internal.NewServiceWithRepositoryAndOwnerVerifier(repository, profileReader)
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	mux := internal.NewMuxWithBus(service, bus)
	addr := ":8086"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("progress-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
