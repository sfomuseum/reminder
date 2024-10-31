package process

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sfomuseum/go-messenger"
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

	db, err := database.NewRemindersDatabase(ctx, opts.RemindersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create reminder database, %w", err)
	}

	defer db.Close()

	m, err := messenger.NewMultiDeliveryAgentWithURIs(ctx, opts.MessengerAgentURIs...)

	if err != nil {
		return fmt.Errorf("Failed to create messenger, %w", err)
	}

	wg := new(sync.WaitGroup)

	process := func(ctx context.Context) error {

		for r, err := range db.Reminders(ctx) {

			if err != nil {
				slog.Error("Failed to retrieve reminder", "error", err)
				break
			}

			wg.Add(1)

			go func(r *reminder.Reminder) {

				defer wg.Done()

				logger := slog.Default()
				logger = logger.With("reminder", r.Id)

				is_due, err := r.IsDue()

				if err != nil {
					slog.Error("Failed to determine is reminder is due", "error", err)
				}

				if !is_due {
					return
				}

				msg := &messenger.Message{
					To:      r.DeliverTo,
					Subject: "Reminder",
					Body:    r.Message,
				}

				err = m.DeliverMessage(ctx, msg)

				if err != nil {
					logger.Error("Failed to deliver reminder", "id", r.Id, "to", r.DeliverTo, "error", err)
				}

			}(r)

		}

		wg.Wait()
		return nil
	}

	switch opts.Mode {
	case "cli":

		err := process(ctx)

		if err != nil {
			logger.Error("Failed to process reminders", "error", err)
			return fmt.Errorf("Failed to process reminders, %w", err)
		}

	case "daemon":

		ticker := time.NewTicker(60 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:

				err := process(ctx)

				if err != nil {
					logger.Error("Failed to process reminders", "error", err)
				}
			}
		}

	case "lambda":

		return fmt.Errorf("Not implemented")

	default:
		return fmt.Errorf("Unsupported mode")
	}

	return nil
}
