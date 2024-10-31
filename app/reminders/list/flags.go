package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var reminders_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("process")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
