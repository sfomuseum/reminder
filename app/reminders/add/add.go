package add

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/adhocore/gronx"
	"github.com/sfomuseum/iso8601duration"
	"github.com/sfomuseum/reminder"
	"github.com/sfomuseum/reminder/database"
)

func Run(ctx context.Context) error {

	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	db, err := database.NewRemindersDatabase(ctx, opts.RemindersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create reminder database, %w", err)
	}

	defer db.Close()

	_, err = gronx.NextTick(opts.Schedule, true)

	if err != nil {
		return fmt.Errorf("Failed to parse schedule, %w", err)
	}

	_, err = duration.FromString(opts.NotifyBefore)

	if err != nil {
		return fmt.Errorf("Failed to parse notify before string, %w", err)
	}

	if opts.Message == "" {
		return fmt.Errorf("Message is empty")
	}

	id, err := reminder.NewId()

	if err != nil {
		return fmt.Errorf("Failed to create new reminder ID, %w", err)
	}

	r := &reminder.Reminder{
		Id:           id,
		Schedule:     opts.Schedule,
		NotifyBefore: opts.NotifyBefore,
		Message:      opts.Message,
		DeliverTo:    opts.DeliverTo,
	}

	err = db.AddReminder(ctx, r)

	if err != nil {
		return fmt.Errorf("Failed to add reminder, %w", err)
	}

	logger.Info("New reminder added", "id", id)
	return nil
}
