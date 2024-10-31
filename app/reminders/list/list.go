package list

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/sfomuseum/go-csvdict"
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
		logger.Debug("Verbose logging enabled")
	}

	db, err := database.NewRemindersDatabase(ctx, opts.RemindersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create reminder database, %w", err)
	}

	defer db.Close()

	var csv_wr *csvdict.Writer

	for r, err := range db.Reminders(ctx) {

		if err != nil {
			slog.Error("Failed to retrieve reminder", "error", err)
			break
		}

		row := map[string]string{
			"id":            strconv.FormatInt(r.Id, 10),
			"schedule":      r.Schedule,
			"notify_before": r.NotifyBefore,
			"message":       r.Message,
			"deliver_to":    r.DeliverTo,
			"deliver_from":  r.DeliverFrom,
		}

		if csv_wr == nil {

			fieldnames := make([]string, 0)

			for k, _ := range row {
				fieldnames = append(fieldnames, k)
			}

			wr, err := csvdict.NewWriter(os.Stdout, fieldnames)

			if err != nil {
				return fmt.Errorf("Failed to create CSV writer, %w", err)
			}

			csv_wr = wr
			csv_wr.WriteHeader()
		}

		err = csv_wr.WriteRow(row)

		if err != nil {
			return fmt.Errorf("Failed to row reminder row (%d), %w", r.Id, err)
		}
	}

	csv_wr.Flush()
	return nil
}
