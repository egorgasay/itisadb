package transactionlogger

import (
	"grpc-storage/internal/grpc-storage/transaction-logger/db"
	"grpc-storage/internal/grpc-storage/transaction-logger/file"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
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

func NewTransactionLogger(Type string, dir string) (ITransactionLogger, error) {
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return nil, err
	}
	switch Type {
	case DB:
		return db.NewLogger(dir, false)
	case VDB:
		return db.NewLogger(dir, true)
	default: // File logger by default
		return file.NewLogger(dir)
	}
}
