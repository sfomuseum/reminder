package process

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/sfomuseum/go-messenger"
	"github.com/sfomuseum/reminder"
	"github.com/sfomuseum/reminder/database"
)

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	db, err := database.NewReminderDatabase(ctx, opts.RemindersDatabaseURI)

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

		for r, err := range db.PendingReminders(ctx) {

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
					Subject: "Reminder",
					Body:    r.Message,
				}

				err = m.DeliverMessage(ctx, msg)

				if err != nil {
					logger.Error("Failed to deliver reminder", "id", r.Id, "error", err)
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

	case "lambda":

		return fmt.Errorf("Not implemented")

	default:
		return fmt.Errorf("Unsupported mode")
	}

	return nil
}
