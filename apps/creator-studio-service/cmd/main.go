package main

import (
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/creator-studio-service/internal"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func main() {
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	service, err := internal.NewService(bus)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := service.Close(); err != nil {
			log.Printf("close creator service outbox: %v", err)
		}
	}()
	mux := internal.NewMux(service)
	addr := ":8087"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("creator-studio-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
