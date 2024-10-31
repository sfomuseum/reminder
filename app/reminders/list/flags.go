package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var reminders_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("process")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "A valid sfomuseum/reminder/database.RemindersDatabase URI.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
