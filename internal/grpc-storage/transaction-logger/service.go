package transactionlogger

import (
	"os"
	"sync"
)

type EventType byte

const (
	Set EventType = iota
	Delete
)

type Event struct {
	EventType EventType
	Key       string
	Value     string
}

type TransactionLogger struct {
	path string
	file *os.File

	events chan Event
	errors chan error

	sync.RWMutex
}

func New(dir string) (*TransactionLogger, error) {
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return nil, err
	}

	tLogger, err := NewLogger(dir)
	if err != nil {
		return nil, err
	}

	return tLogger, nil
}
