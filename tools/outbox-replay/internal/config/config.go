package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"
)

type Mode string

const (
	ModeListFailed    Mode = "list-failed"
	ModeRequeueFailed Mode = "requeue-failed"
	ModeRequeueEvent  Mode = "requeue-event"
)

type Options struct {
	DatabaseURL   string
	Mode          Mode
	Limit         int
	Topic         string
	EventID       string
	DryRun        bool
	ResetAttempts bool
	Timeout       time.Duration
}

func Parse(args []string, getenv func(string) string) (Options, error) {
	var opts Options
	var modeRaw string
	fs := flag.NewFlagSet("outbox-replay", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.DatabaseURL, "database-url", "", "Postgres database URL (falls back to DATABASE_URL)")
	fs.StringVar(&modeRaw, "mode", string(ModeListFailed), "Mode: list-failed|requeue-failed|requeue-event")
	fs.IntVar(&opts.Limit, "limit", 25, "Max rows for list-failed/requeue-failed")
	fs.StringVar(&opts.Topic, "topic", "", "Optional topic filter for list-failed/requeue-failed")
	fs.StringVar(&opts.EventID, "event-id", "", "Event ID for requeue-event mode")
	fs.BoolVar(&opts.DryRun, "dry-run", true, "Preview-only for requeue modes")
	fs.BoolVar(&opts.ResetAttempts, "reset-attempts", true, "Reset attempts to 0 when requeueing")
	fs.DurationVar(&opts.Timeout, "timeout", 15*time.Second, "Overall command timeout")
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	if opts.DatabaseURL == "" {
		opts.DatabaseURL = getenv("DATABASE_URL")
	}
	mode, err := parseMode(modeRaw)
	if err != nil {
		return Options{}, err
	}
	opts.Mode = mode
	if err := opts.validate(); err != nil {
		return Options{}, err
	}
	return opts, nil
}

func parseMode(raw string) (Mode, error) {
	switch Mode(raw) {
	case ModeListFailed, ModeRequeueFailed, ModeRequeueEvent:
		return Mode(raw), nil
	default:
		return "", fmt.Errorf("unsupported mode %q", raw)
	}
}

func (o Options) validate() error {
	if o.DatabaseURL == "" {
		return errors.New("database URL is required (flag -database-url or env DATABASE_URL)")
	}
	if o.Limit <= 0 {
		return errors.New("limit must be > 0")
	}
	if o.Timeout <= 0 {
		return errors.New("timeout must be > 0")
	}
	if o.Mode == ModeRequeueEvent && o.EventID == "" {
		return errors.New("event-id is required in requeue-event mode")
	}
	return nil
}
