package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/delqhi/mikasmissions/platform/workers/worker-social-snippet/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	processor := internal.NewProcessor("worker-social-snippet", "episode.published.v1", logger)
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := bus.Subscribe(ctx, "episode.published.v1", "worker-social-snippet", processor.Handle); err != nil {
		log.Fatal(err)
	}
	logger.Info("worker started", "worker", "worker-social-snippet", "topic", "episode.published.v1")
	<-ctx.Done()
}
