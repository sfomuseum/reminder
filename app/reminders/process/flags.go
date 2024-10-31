package process

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var reminders_database_uri string
var messenger_agents_uris multi.MultiString

var mode string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("process")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "...")
	fs.Var(&messenger_agents_uris, "messenger-agent-uri", "...")
	fs.StringVar(&mode, "mode", "cli", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
