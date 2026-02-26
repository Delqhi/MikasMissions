package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/delqhi/mikasmissions/platform/workers/worker-gen-orchestrator/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bus, err := queue.NewBusFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	processor := internal.NewProcessor(bus, logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := bus.Subscribe(ctx, processor.Topic(), processor.Consumer(), processor.Handle); err != nil {
		log.Fatal(err)
	}
	logger.Info("worker started", "worker", processor.Consumer(), "topic", processor.Topic())
	<-ctx.Done()
}
