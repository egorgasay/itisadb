package transactionlogger

import (
	"os"
	"sync"
)

type EventType byte

const (
	Set EventType = iota
	Delete
	SetToIndex
	DeleteAttr
	CreateIndex
	Attach
	DeleteIndex
)

type Event struct {
	EventType EventType
	Name      string
	Value     string
}

type TransactionLogger struct {
	path string

	events chan Event
	errors chan error

	sync.RWMutex
}

func New() (*TransactionLogger, error) {
	if err := os.MkdirAll(PATH, 0755); err != nil {
		return nil, err
	}

	f, err := os.Create(PATH + "/1")
	if err != nil {
		return nil, err
	}
	f.Close()

	return &TransactionLogger{
		path: PATH,
	}, nil
}
