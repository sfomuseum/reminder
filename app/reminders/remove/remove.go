package remove

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

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

	for _, id := range opts.Ids {

		r := &reminder.Reminder{
			Id: id,
		}

		err = db.RemoveReminder(ctx, r)

		if err != nil {
			return fmt.Errorf("Failed to remove reminder %d, %w", id, err)
		}

		logger.Info("Reminder has been removed", "id", id)
	}

	return nil
}
