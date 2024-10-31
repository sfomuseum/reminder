package reminder

import (
	"time"

	"github.com/adhocore/gronx"
	"github.com/sfomuseum/iso8601duration"
)

type Reminder struct {
	Id        int64  `json:"id"`
	Schedule  string `json:"schedule"`
	// An ISO8601 duration string
	NotifyBefore string `json:"notify_before"`
	Message      string `json:"message"`
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
