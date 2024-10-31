package database

import (
	"context"
	"iter"

	"github.com/sfomuseum/reminder"
)

type RemindersDatabase interface {
	AddReminder(context.Context, *reminder.Reminder) error
	UpdateReminder(context.Context, *reminder.Reminder) error
	PendingReminders(context.Context) iter.Seq2[*reminder.Reminder, error]
	Close() error
}
