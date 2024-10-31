package database

import (
	"context"
	"fmt"
	"io"
	"iter"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/reminder"
)

type CSVRemindersDatabase struct {
	RemindersDatabase
	path string
}

func init() {
	ctx := context.Background()
	RegisterRemindersDatabase(ctx, "csv", NewCSVRemindersDatabase)
}

func NewCSVRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(u.Path)

	if err != nil {
		return nil, err
	}

	db := &CSVRemindersDatabase{
		path: path,
	}

	return db, nil
}

func (db *CSVRemindersDatabase) Close() error {
	return nil
}

func (db *CSVRemindersDatabase) AddReminder(ctx context.Context, r *reminder.Reminder) error {
	return fmt.Errorf("Not implemented")
}

func (db *CSVRemindersDatabase) RemoveReminder(ctx context.Context, r *reminder.Reminder) error {
	return fmt.Errorf("Not implemented")
}

func (db *CSVRemindersDatabase) Reminders(ctx context.Context) iter.Seq2[*reminder.Reminder, error] {

	return func(yield func(*reminder.Reminder, error) bool) {

		csv_r, err := csvdict.NewReaderFromPath(db.path)

		if err != nil {
			yield(nil, err)
			return
		}

		for {
			row, err := csv_r.Read()

			if err == io.EOF {
				return
			}

			if err != nil {
				yield(nil, err)
				return
			}

			id, err := strconv.ParseInt(row["id"], 10, 64)

			if err != nil {
				yield(nil, err)
				continue
			}

			r := &reminder.Reminder{
				Id:           id,
				Schedule:     row["schedule"],
				NotifyBefore: row["notify_before"],
				Message:      row["message"],
				DeliverTo:    row["deliver_to"],
			}

			yield(r, nil)
		}

	}
}
