package remove

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var reminders_database_uri string
var ids multi.MultiInt64

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("add")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "A valid sfomuseum/reminder/database.RemindersDatabase URI.")
	fs.Var(&ids, "id", "One or more reminder IDs to remove")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
