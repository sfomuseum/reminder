package database

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"net/url"

	"github.com/sfomuseum/reminder"
)

type SQLRemindersDatabase struct {
	RemindersDatabase
	conn *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterRemindersDatabase(ctx, "sql", NewSQLRemindersDatabase)
}

func NewSQLRemindersDatabase(ctx context.Context, uri string) (RemindersDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	engine := u.Host

	q := u.Query()

	dsn := q.Get("dsn")

	conn, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	db := &SQLRemindersDatabase{
		conn: conn,
	}

	return db, nil
}

func (db *SQLRemindersDatabase) Close() error {
	return db.conn.Close()
}

func (db *SQLRemindersDatabase) AddReminder(ctx context.Context, r *reminder.Reminder) error {

	q := "INSERT INTO reminders (id, schedule, notify_before, message, deliver_to) VALUES (?, ?, ?, ?, ?)"
	_, err := db.conn.ExecContext(ctx, q, r.Id, r.Schedule, r.NotifyBefore, r.Message, r.DeliverTo)
	return err
}

func (db *SQLRemindersDatabase) RemoveReminder(ctx context.Context, r *reminder.Reminder) error {

	q := "DELETE FROM reminders WHERE id = ?"
	_, err := db.conn.ExecContext(ctx, q, r.Id)
	return err
}

func (db *SQLRemindersDatabase) Reminders(ctx context.Context) iter.Seq2[*reminder.Reminder, error] {

	return func(yield func(*reminder.Reminder, error) bool) {

		q := "SELECT id, schedule, notify_before, message, deliver_to FROM reminders"

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var id int64
			var schedule string
			var notify_before string
			var message string
			var deliver_to string

			err := rows.Scan(&id, &schedule, &notify_before, &message, &deliver_to)

			if err != nil {
				yield(nil, err)
				return
			}

			r := &reminder.Reminder{
				Id:           id,
				Schedule:     schedule,
				NotifyBefore: notify_before,
				Message:      message,
				DeliverTo:    deliver_to,
			}

			yield(r, nil)
		}

		err = rows.Err()

		if err != nil {
			yield(nil, err)
		}

	}
}
