package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/delqhi/mikasmissions/platform/tools/outbox-replay/internal/config"
	"github.com/delqhi/mikasmissions/platform/tools/outbox-replay/internal/store"
)

func main() {
	opts, err := config.Parse(os.Args[1:], os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n\n", err)
		printUsage()
		os.Exit(2)
	}
	repo, err := store.New(opts.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open outbox store: %v\n", err)
		os.Exit(1)
	}
	defer repo.Close()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	switch opts.Mode {
	case config.ModeListFailed:
		runListFailed(ctx, repo, opts)
	case config.ModeRequeueFailed:
		runRequeueFailed(ctx, repo, opts)
	case config.ModeRequeueEvent:
		runRequeueEvent(ctx, repo, opts)
	default:
		fmt.Fprintf(os.Stderr, "unsupported mode %q\n", opts.Mode)
		os.Exit(2)
	}
}

func runListFailed(ctx context.Context, repo *store.Store, opts config.Options) {
	rows, err := repo.ListFailed(ctx, opts.Topic, opts.Limit)
	if err != nil {
		exitErr("list failed rows", err)
	}
	if len(rows) == 0 {
		fmt.Println("no failed outbox rows")
		return
	}
	fmt.Println("event_id\ttopic\tattempts\tavailable_at\tlast_error")
	for _, row := range rows {
		fmt.Printf(
			"%s\t%s\t%d\t%s\t%s\n",
			row.EventID,
			row.Topic,
			row.Attempts,
			row.AvailableAt.UTC().Format(time.RFC3339),
			row.LastError,
		)
	}
}

func runRequeueFailed(ctx context.Context, repo *store.Store, opts config.Options) {
	rows, err := repo.RequeueFailed(ctx, opts.Topic, opts.Limit, opts.DryRun, opts.ResetAttempts)
	if err != nil {
		exitErr("requeue failed rows", err)
	}
	printReplayResult(rows, opts.DryRun)
}

func runRequeueEvent(ctx context.Context, repo *store.Store, opts config.Options) {
	rows, err := repo.RequeueEvent(ctx, opts.EventID, opts.DryRun, opts.ResetAttempts)
	if err != nil {
		exitErr("requeue event", err)
	}
	printReplayResult(rows, opts.DryRun)
}

func printReplayResult(eventIDs []string, dryRun bool) {
	modeLabel := "requeued"
	if dryRun {
		modeLabel = "candidates"
	}
	if len(eventIDs) == 0 {
		fmt.Printf("no %s found\n", modeLabel)
		return
	}
	fmt.Printf("%s (%d):\n", modeLabel, len(eventIDs))
	for _, eventID := range eventIDs {
		fmt.Printf("- %s\n", eventID)
	}
}

func exitErr(op string, err error) {
	fmt.Fprintf(os.Stderr, "%s failed: %v\n", op, err)
	os.Exit(1)
}

func printUsage() {
	fmt.Println("Usage: go run ./tools/outbox-replay/cmd -mode=<mode> [flags]")
	fmt.Println("")
	fmt.Println("Modes:")
	fmt.Println("  list-failed      List failed outbox rows")
	fmt.Println("  requeue-failed   Requeue failed rows (supports dry-run)")
	fmt.Println("  requeue-event    Requeue one failed row by event-id")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run ./tools/outbox-replay/cmd -mode=list-failed -limit=20")
	fmt.Println("  go run ./tools/outbox-replay/cmd -mode=requeue-failed -topic=episode.published.v1 -dry-run=true")
	fmt.Println("  go run ./tools/outbox-replay/cmd -mode=requeue-event -event-id=evt-123 -dry-run=false")
}
