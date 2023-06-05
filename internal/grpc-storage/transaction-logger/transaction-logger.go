package transactionlogger

import (
	"itisadb/internal/grpc-storage/transaction-logger/db"
	"itisadb/internal/grpc-storage/transaction-logger/file"
	"itisadb/internal/grpc-storage/transaction-logger/service"
	"os"
)

type ITransactionLogger interface {
	WriteSet(key, value string)
	WriteDelete(key string)
	Err() <-chan error

	ReadEvents() (<-chan service.Event, <-chan error)

	Run()
	Clear() error
}

const File = "file"
const DB = "db"
const VDB = "docker_db"

// NewTransactionLogger creates new transaction logger
// can return nil, nil
func NewTransactionLogger(Type string, creds string) (ITransactionLogger, error) {
	switch Type {
	case DB:
		return db.NewLogger(creds, false)
	case VDB:
		return db.NewLogger(creds, true)
	case File:
		err := os.MkdirAll(creds, 0644)
		if err != nil {
			return nil, err
		}
		return file.NewLogger(creds)
	}
	return nil, nil
}
