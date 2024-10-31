package remove

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	RemindersDatabaseURI string
	Ids                  []int64
	Verbose              bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		RemindersDatabaseURI: reminders_database_uri,
		Ids:                  ids,
		Verbose:              verbose,
	}

	return opts, nil
}
