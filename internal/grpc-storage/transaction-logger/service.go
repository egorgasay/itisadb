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
	pathToDir  string
	pathToFile string
	file       *os.File

	currentCOL  int32
	currentName int32

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

	f, err = os.OpenFile(PATH+"/1", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &TransactionLogger{
		pathToDir:   PATH,
		pathToFile:  PATH + "/1",
		file:        f,
		currentName: 1,
	}, nil
}
