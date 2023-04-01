package transactionlogger

import (
	"grpc-storage/internal/grpc-storage/transaction-logger/file"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
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

func NewTransactionLogger(Type string) (ITransactionLogger, error) {
	switch Type {
	default:
		return file.NewLogger()
	}
}
