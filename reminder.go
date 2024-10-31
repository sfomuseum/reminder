package reminder

import (
	"log/slog"
	"time"

	"github.com/adhocore/gronx"
	"github.com/sfomuseum/iso8601duration"
)

// Reminder is a struct encapsulating details about a reminder message.
type Reminder struct {
	// A unique identifier for the reminder.
	Id int64 `json:"id"`
	// A valid cron expression (that can be parsed by adhocore/gronx) for the scheduled event.
	Schedule string `json:"schedule"`
	// An ISO8601 duration string indicating the amount of time before the scheduled event is due to start sending reminders.
	NotifyBefore string `json:"notify_before"`
	// The message body of the reminder.
	Message string `json:"message"`
	// The address where the reminder should be delivered to.
	DeliverTo string `json:deliver_to"`
	// The address where the reminder should be delivered from.
	DeliverFrom string `json:"deliver_from"`
}

// Return a boolean flag where the reminder is due to be dispatched.
func (r *Reminder) IsDue() (bool, error) {

	next, err := gronx.NextTick(r.Schedule, true)

	if err != nil {
		return false, err
	}

	logger := slog.Default()
	logger = logger.With("reminder", r.Id)

	dur, err := duration.FromString(r.NotifyBefore)

	if err != nil {
		return false, err
	}

	d := dur.ToDuration()

	now := time.Now()

	logger.Debug("Now", "t", now)
	logger.Debug("Next", "t", next)

	trigger := next.Add(-d)

	logger.Debug("Trigger", "t", trigger)

	if trigger.After(now) {
		return false, nil
	}

	logger.Debug("Reminder is due")
	return true, nil
}
