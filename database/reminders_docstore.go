package database

import (
	"context"
	"iter"

	"github.com/sfomuseum/reminder"
	"gocloud.dev/docstore"
)

type DocstoreRemindersDatabase struct {
	RemindersDatabase
	collection *docstore.Collection
}

func NewDocstoreRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	db := &DocstoreRemindersDatabase{}

	return db, nil
}

func (db *DocstoreRemindersDatabase) AddReminder(ctx context.Context, r *reminder.Reminder) error {
	return nil
}

func (db *DocstoreRemindersDatabase) UpdateReminder(ctx context.Context, r *reminder.Reminder) error {
	return nil
}

func (db *DocstoreRemindersDatabase) PendingReminders(ctx context.Context) iter.Seq2[*reminder.Reminder, error] {
	return nil
}
