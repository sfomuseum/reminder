package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var reminders_database_uri string

var schedule string
var notify_before string
var message string
var deliver_to string
var deliver_from string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("add")

	fs.StringVar(&reminders_database_uri, "reminders-database-uri", "", "A valid sfomuseum/reminder/database.RemindersDatabase URI.")
	fs.StringVar(&schedule, "schedule", "", "A valid cron expression (that can be parsed by adhocore/gronx) for the scheduled event.")
	fs.StringVar(&notify_before, "notify-before", "", "An ISO8601 duration string indicating the amount of time before the scheduled event is due to start sending reminders.")
	fs.StringVar(&message, "message", "", "The message body of the reminder.")
	fs.StringVar(&deliver_to, "deliver-to", "", "The address where the reminder should be delivered to.")
	fs.StringVar(&deliver_from, "deliver-from", "", "The address where the reminder should be delivered from.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
