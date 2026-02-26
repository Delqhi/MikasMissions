package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/delqhi/mikasmissions/platform/workers/worker-outbox-relay/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required for worker-outbox-relay")
	}
	outbox, err := queue.NewPersistentOutbox(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer outbox.Close()

	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	interval := readRelayInterval()
	relay := internal.NewRelay(outbox, bus, interval, logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Info("worker started", "worker", "worker-outbox-relay", "interval_ms", interval.Milliseconds())
	relay.Run(ctx)
}

func readRelayInterval() time.Duration {
	raw := os.Getenv("OUTBOX_RELAY_INTERVAL_MS")
	if raw == "" {
		return 1000 * time.Millisecond
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 1000 * time.Millisecond
	}
	return time.Duration(parsed) * time.Millisecond
}
