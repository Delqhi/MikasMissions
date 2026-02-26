package main

import (
	"log"
	"net/http"
	"os"

	"github.com/delqhi/mikasmissions/platform/apps/playback-service/internal"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func main() {
	profileReader, err := profileclient.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	service := internal.NewServiceWithDependencies(
		internal.NewBillingEntitlementVerifierFromEnv(),
		internal.NewIdentityGateVerifierFromEnv(),
		profileReader,
	)
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	mux := internal.NewMuxWithBus(service, bus)
	addr := ":8085"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		addr = ":" + fromEnv
	}
	log.Printf("playback-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
