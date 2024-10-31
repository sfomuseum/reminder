package process

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var reminders_database_uri string
var messenger_agents_uris multi.MultiCSVString

var mode string
var verbose bool

var frequency string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("process")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "A valid sfomuseum/reminder/database.RemindersDatabase URI.")
	fs.Var(&messenger_agents_uris, "messenger-agent-uri", "One or more valid sfomuseum/go-messenger.Messenger URIs.")
	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, daemon, lambda")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.StringVar(&frequency, "frequency", "PT1M", "A valid ISO8601 duration string indicating how often to process reminders. Required if -mode daemon.")
	return fs
}
