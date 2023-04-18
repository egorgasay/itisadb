package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/egorgasay/dockerdb/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"log"
	"strings"
)

type TransactionLogger struct {
	db   *sql.DB
	path string

	events chan service.Event
	errors chan error
}

const insertQuery = "INSERT INTO transactions (`event_type`, `wrench`, `value`) VALUES " +
	" (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?)," +
	"(?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), " +
	" (?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), " +
	"(?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?), " +
	"(?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?)"

var insertEvent *sql.Stmt

func NewLogger(path string, vdb bool) (*TransactionLogger, error) {
	path = strings.TrimRight(path, "/") + "/transactionLoggerDB"

	var db *sql.DB
	if vdb {
		cfg := dockerdb.CustomDB{
			DB: dockerdb.DB{
				Name:     "a34dm8",
				User:     "adm",
				Password: "adm",
			},
			Port:   "3247",
			Vendor: dockerdb.MySQL8Image,
		}

		dockerDB, err := dockerdb.New(context.Background(), cfg)
		if err != nil {
			return nil, err
		}

		db = dockerDB.DB
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/mysql_tlogger",
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

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	insertEvent, err = db.Prepare(insertQuery)
	if err != nil {
		return nil, err
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
	events := make(chan service.Event, 30000)

	t.events = events
	var data = make([]service.Event, 0, 30000)
	var dataBackup = make([]service.Event, 30000)
	go func() {
		for e := range events {
			data = append(data, e)
			if len(data) == 20 {
				copy(dataBackup, data)
				go t.flash(dataBackup)
				data = data[:0]
			}
		}
	}()
}

// flash grabs 20 events and saves them to the db.
func (t *TransactionLogger) flash(data []service.Event) {
	errorsch := make(chan error, 20)
	t.errors = errorsch

	var anys = make([]any, 0, len(data)*3)

	for _, e := range data {
		anys = append(anys, e.EventType, e.Key, e.Value)
	}

	_, err := insertEvent.Exec(anys...)
	if err != nil {
		log.Println("Run:", err)
		errorsch <- err
	}
}

func (t *TransactionLogger) Err() <-chan error {
	return t.errors
}

func (t *TransactionLogger) ReadEvents() (<-chan service.Event, <-chan error) {
	outEvent := make(chan service.Event)
	outError := make(chan error, 1)

	rows, err := t.db.Query("SELECT `event_type`, `wrench`, `value` FROM transactions")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		outError <- err
		return outEvent, outError
	}
	defer rows.Close()

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
