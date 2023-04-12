package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TransactionLogger struct {
	db   *sql.DB
	path string

	events chan service.Event
	errors chan error
}

func NewLogger(path string) (*TransactionLogger, error) {
	path += "/transactionLoggerDB"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/sqlite3_tlogger",
		"sqlite", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatal(err)
		}
	}

	return &TransactionLogger{
		db:   db,
		path: path,
	}, nil
}

func (t *TransactionLogger) WriteSet(key, value string) {
	t.events <- service.Event{EventType: service.Set, Key: key, Value: value}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- service.Event{EventType: service.Delete, Key: key}
}

func (t *TransactionLogger) Run() {
	events := make(chan service.Event, 20)
	errorsch := make(chan error, 20)

	t.events = events
	t.errors = errorsch

	go func() {
		for e := range events {
			_, err := t.db.Exec("INSERT INTO transactions (event_type, key, value) VALUES (?, ?, ?)", e.EventType, e.Key, e.Value)
			if err != nil {
				log.Println("Run:", err)
				errorsch <- err
				return
			}
		}
	}()
}
func (t *TransactionLogger) Err() <-chan error {
	return t.errors
}

func (t *TransactionLogger) ReadEvents() (<-chan service.Event, <-chan error) {
	outEvent := make(chan service.Event)
	outError := make(chan error, 1)

	rows, err := t.db.Query("SELECT event_type, key, value FROM transactions")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		outError <- err
		return outEvent, outError
	}

	go func() {
		var event service.Event
		defer close(outEvent)
		defer close(outError)

		for rows.Next() {
			err = rows.Scan(&event.EventType, &event.Key, &event.Value)
			if err != nil {
				log.Println("ReadEvents:", err)
				continue
			}

			outEvent <- event
		}
	}()

	return outEvent, outError
}

func (t *TransactionLogger) Clear() error {
	_, err := t.db.Exec("DELETE FROM transactions")
	if err != nil {
		return err
	}

	return nil
}
