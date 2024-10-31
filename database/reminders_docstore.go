package database

import (
	"context"
	"fmt"
	"io"
	"iter"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/reminder"
	"gocloud.dev/docstore"
)

type DocstoreRemindersDatabase struct {
	RemindersDatabase
	collection *docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterRemindersDatabase(ctx, "awsdynamodb", NewDocstoreRemindersDatabase)

	for _, scheme := range docstore.DefaultURLMux().CollectionSchemes() {
		RegisterRemindersDatabase(ctx, scheme, NewDocstoreRemindersDatabase)
	}
}

func NewDocstoreRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreRemindersDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreRemindersDatabase) Close() error {
	return db.collection.Close()
}

func (db *DocstoreRemindersDatabase) AddReminder(ctx context.Context, r *reminder.Reminder) error {
	return db.collection.Put(ctx, r)
}

func (db *DocstoreRemindersDatabase) RemoveReminder(ctx context.Context, r *reminder.Reminder) error {
	return db.collection.Delete(ctx, r)
}

func (db *DocstoreRemindersDatabase) Reminders(ctx context.Context) iter.Seq2[*reminder.Reminder, error] {
	q := db.collection.Query()
	return db.getRemindersWithQuery(ctx, q)
}

func (db *DocstoreRemindersDatabase) getRemindersWithQuery(ctx context.Context, q *docstore.Query) iter.Seq2[*reminder.Reminder, error] {

	return func(yield func(*reminder.Reminder, error) bool) {

		iter := q.Get(ctx)
		defer iter.Stop()

		for {

			var r reminder.Reminder
			err := iter.Next(ctx, &r)

			if err == io.EOF {
				break
			} else if err != nil {
				yield(nil, err)
			} else {
				yield(&r, nil)
			}
		}
	}
}
