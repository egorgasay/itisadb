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
	kvPath      string
	indexesPath string

	events chan Event
	errors chan error

	sync.RWMutex
}

func New() (*TransactionLogger, error) {
	if err := os.MkdirAll(PATH+"/kv", 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(PATH+"/indexes", 0755); err != nil {
		return nil, err
	}

	return &TransactionLogger{
		indexesPath: PATH + "/indexes",
		kvPath:      PATH + "/kv",
	}, nil
}
