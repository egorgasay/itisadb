package transactionlogger

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type EventType byte

const (
	Set EventType = iota
	Delete
	SetToObject
	DeleteAttr
	CreateObject
	Attach
	DeleteObject
	CreateUser
	DeleteUser
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

	files, err := os.ReadDir(PATH)
	if err != nil {
		return nil, err
	}

	maxNumber := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if n, err := strconv.Atoi(f.Name()); err != nil {
			continue
		} else if n > maxNumber {
			maxNumber = n
		}
	}

	if maxNumber == 0 {
		maxNumber = 1
	}

	filename := fmt.Sprint(PATH, "/", maxNumber)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &TransactionLogger{
		pathToDir:   PATH,
		pathToFile:  filename,
		file:        f,
		currentName: int32(maxNumber),
	}, nil
}
