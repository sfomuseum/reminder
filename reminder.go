package reminder

import (
	"time"

	"github.com/adhocore/gronx"
	"github.com/sfomuseum/iso8601duration"
)

type Reminder struct {
	// A unique identifier for the reminder
	Id int64 `json:"id"`
	// A valid cron expression (that can be parsed by adhocore/gronx)
	Schedule string `json:"schedule"`
	// An ISO8601 duration string
	NotifyBefore string `json:"notify_before"`
	// The message body of the reminder
	Message string `json:"message"`
	// A list of addresses to deliver the reminder to
	DeliverTo string `json:deliver_to"`
}

func (r *Reminder) IsDue() (bool, error) {

	next, err := gronx.NextTick(r.Schedule, true)

	if err != nil {
		return false, err
	}

	dur, err := duration.FromString(r.NotifyBefore)

	if err != nil {
		return false, err
	}

	d := dur.ToDuration()

	trigger := next.Add(-d)

	now := time.Now()

	if trigger.After(now) {
		return false, nil
	}

	return true, nil
}
