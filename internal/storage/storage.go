package storage

import (
	"database/sql"
	"errors"
	"github.com/egorgasay/grpc-storage/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

type Storage struct {
	DBStore    *sql.DB
	RAMStorage map[string]string
}

func New(cfg *config.DBConfig) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty configuration")
	}

	db, err := sql.Open(cfg.DriverName, cfg.DataSourceCred)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	err = m.Up()
	if err.Error() != "no change" {
		log.Fatal(err)
	}

	return &Storage{
		DBStore:    db,
		RAMStorage: make(map[string]string, 10),
	}, nil
}
