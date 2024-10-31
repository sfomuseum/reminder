package database

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/reminder"
)

type RemindersDatabase interface {
	AddReminder(context.Context, *reminder.Reminder) error
	RemoveReminder(context.Context, *reminder.Reminder) error
	Reminders(context.Context) iter.Seq2[*reminder.Reminder, error]
	Close() error
}

var reminders_database_roster roster.Roster

// RemindersDatabaseInitializationFunc is a function defined by individual reminders_database package and used to create
// an instance of that reminders_database
type RemindersDatabaseInitializationFunc func(ctx context.Context, uri string) (RemindersDatabase, error)

// RegisterRemindersDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `RemindersDatabase` instances by the `NewRemindersDatabase` method.
func RegisterRemindersDatabase(ctx context.Context, scheme string, init_func RemindersDatabaseInitializationFunc) error {

	err := ensureRemindersDatabaseRoster()

	if err != nil {
		return err
	}

	return reminders_database_roster.Register(ctx, scheme, init_func)
}

func ensureRemindersDatabaseRoster() error {

	if reminders_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		reminders_database_roster = r
	}

	return nil
}

// NewRemindersDatabase returns a new `RemindersDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `RemindersDatabaseInitializationFunc`
// function used to instantiate the new `RemindersDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterRemindersDatabase` method.
func NewRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := reminders_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(RemindersDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func RemindersDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureRemindersDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range reminders_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
