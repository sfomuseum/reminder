package process

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sfomuseum/go-messenger"
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
		logger.Debug("Verbose logging enabled")
	}

	messenger.RegisterEmailSchemes(ctx)

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
				logger.Error("Failed to retrieve reminder", "error", err)
				break
			}

			wg.Add(1)

			go func(r *reminder.Reminder) {

				defer wg.Done()

				logger := slog.Default()
				logger = logger.With("reminder", r.Id)

				logger.Debug("Check whether reminder is due")
				
				is_due, err := r.IsDue()

				if err != nil {
					logger.Error("Failed to determine if reminder is due", "error", err)
				}

				if !is_due {
					logger.Debug("Reminder is not due, skipping")
					return
				}

				subject := fmt.Sprintf("Reminder #%d", r.Id)

				msg := &messenger.Message{
					To:      r.DeliverTo,
					From:    r.DeliverFrom,
					Subject: subject,
					Body:    r.Message,
				}

				logger.Debug("Deliver message", "to", r.DeliverTo)
				
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

		d, err := duration.FromString(opts.Frequency)

		if err != nil {
			return fmt.Errorf("Invalid frequency string, %w", err)
		}

		ticker := time.NewTicker(d.ToDuration())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:

				logger.Debug("Process reminders")

				err := process(ctx)

				if err != nil {
					logger.Error("Failed to process reminders", "error", err)
				}
			}
		}

	case "lambda":

		lambda.Start(process)

	default:
		return fmt.Errorf("Unsupported mode")
	}

	return nil
}
