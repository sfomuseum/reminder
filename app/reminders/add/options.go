package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	RemindersDatabaseURI string
	Schedule             string
	NotifyBefore         string
	Message              string
	DeliverTo            string
	Verbose              bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		RemindersDatabaseURI: reminders_database_uri,
		Schedule:             schedule,
		NotifyBefore:         notify_before,
		Message:              message,
		DeliverTo:            deliver_to,
		Verbose:              verbose,
	}

	return opts, nil
}
