package transactionlogger

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"itisadb/config"
)

type EventType byte

const (
	_ EventType = iota
	Set
	Delete
	SetToObject
	DeleteAttr
	CreateObject
	Attach
	DeleteObject
	CreateUser
	DeleteUser
	AddObjectInfo
	DeleteObjectInfo
)

type Event struct {
	EventType EventType
	Name      string
	Value     string
	Metadata  string
}

type TransactionLogger struct {
	pathToFile string
	file       *os.File

	currentCOL  int32
	currentName int32

	events chan Event
	errors chan error

	sync.RWMutex

	logger *zap.Logger
	cfg    config.TransactionLoggerConfig
}

func New(cfg config.TransactionLoggerConfig) (*TransactionLogger, error) {
	if cfg.BackupDirectory == "" {
		cfg.BackupDirectory = DefaultPath
	}

	if err := os.MkdirAll(cfg.BackupDirectory, 0755); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(cfg.BackupDirectory)
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

	filename := fmt.Sprint(cfg.BackupDirectory, "/", maxNumber)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &TransactionLogger{
		pathToFile:  filename,
		file:        f,
		currentName: int32(maxNumber),
		cfg:         cfg,
	}, nil
}
