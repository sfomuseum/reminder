package database

import (
	"context"
	"iter"

	"github.com/sfomuseum/reminder"
)

type NullRemindersDatabase struct {
	RemindersDatabase
	path string
}

func init() {
	ctx := context.Background()
	RegisterRemindersDatabase(ctx, "null", NewNullRemindersDatabase)
}

func NewNullRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	db := &NullRemindersDatabase{}
	return db, nil
}

func (db *NullRemindersDatabase) Close() error {
	return nil
}

func (db *NullRemindersDatabase) AddReminder(ctx context.Context, r *reminder.Reminder) error {
	return nil
}

func (db *NullRemindersDatabase) RemoveReminder(ctx context.Context, r *reminder.Reminder) error {
	return nil
}

func (db *NullRemindersDatabase) Reminders(ctx context.Context) iter.Seq2[*reminder.Reminder, error] {

	return func(yield func(*reminder.Reminder, error) bool) {
		return
	}
}
