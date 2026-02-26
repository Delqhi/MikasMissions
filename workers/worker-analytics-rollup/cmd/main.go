package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/delqhi/mikasmissions/platform/workers/worker-analytics-rollup/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	processor := internal.NewProcessor("worker-analytics-rollup", "analytics.rollup.v1", logger)
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := bus.Subscribe(ctx, "analytics.rollup.v1", "worker-analytics-rollup", processor.Handle); err != nil {
		log.Fatal(err)
	}
	logger.Info("worker started", "worker", "worker-analytics-rollup", "topic", "analytics.rollup.v1")
	<-ctx.Done()
}
