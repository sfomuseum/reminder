package process

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	RemindersDatabaseURI string
	MessengerAgentURIs   []string
	Mode                 string
	Verbose              bool
	Frequency            string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "REMINDER")

	if err != nil {
		return nil, err
	}

	opts := &RunOptions{
		RemindersDatabaseURI: reminders_database_uri,
		MessengerAgentURIs:   messenger_agents_uris,
		Mode:                 mode,
		Verbose:              verbose,
		Frequency:            frequency,
	}

	return opts, nil
}
