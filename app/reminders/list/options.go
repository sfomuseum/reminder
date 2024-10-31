package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	RemindersDatabaseURI string
	Verbose              bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		RemindersDatabaseURI: reminders_database_uri,
		Verbose:              verbose,
	}

	return opts, nil
}
