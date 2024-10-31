package add

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"regexp"
	"time"
	
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

	re_ymd, err := regexp.Compile(`^\d{4}\-\d{2}\-\d{2}$`)

	if err != nil {
		return fmt.Errorf("Failed to compile YMD pattern, %w", err)
	}
		
	schedule := opts.Schedule

	if re_ymd.MatchString(schedule) {

		t, err := time.Parse("2006-01-02", schedule)

		if err != nil {
			return fmt.Errorf("Failed to parse YMD schedule, %w", err)
		}

		schedule = fmt.Sprintf("0 0 %d %d * %d", t.Day(), t.Month(), t.Year())
		logger.Debug("Reassign schedule", "schedule", schedule)
	}
	
	_, err = gronx.NextTick(schedule, true)

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
		DeliverFrom:  opts.DeliverFrom,
	}

	err = db.AddReminder(ctx, r)

	if err != nil {
		return fmt.Errorf("Failed to add reminder, %w", err)
	}

	logger.Info("New reminder added", "id", id)
	return nil
}
