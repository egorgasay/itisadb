package file

import (
	"bufio"
	"fmt"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"strconv"
	"strings"
	"sync"
)

import "os"

type TransactionLogger struct {
	path string
	file *os.File

	events chan service.Event
	errors chan error

	sync.RWMutex
}

func NewLogger(path string) (*TransactionLogger, error) {
	// TODO: ADD CHANGING SERVER NUMBER
	path += "/transactionLogger"

	return &TransactionLogger{
		path: path,
	}, nil
}

func (t *TransactionLogger) WriteSet(key, value string) {
	t.events <- service.Event{EventType: service.Set, Key: key, Value: value}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- service.Event{EventType: service.Delete, Key: key}
}

func (t *TransactionLogger) Run() {
	events := make(chan service.Event, 20)
	errors := make(chan error, 20)

	t.events = events
	t.errors = errors

	go func() {
		t.Lock()
		for e := range events {
			_, err := t.file.WriteString(fmt.Sprintf("%v %s %s\n", e.EventType, e.Key, e.Value))
			if err != nil {
				errors <- err
				return
			}
		}
		t.Unlock()
	}()
}
func (t *TransactionLogger) Err() <-chan error {
	return t.errors
}

func (t *TransactionLogger) ReadEvents() (<-chan service.Event, <-chan error) {
	outEvent := make(chan service.Event)
	outError := make(chan error, 1)
	f, err := os.OpenFile(t.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	t.file = f
	if err != nil {
		outError <- err
		return outEvent, outError
	}
	scanner := bufio.NewScanner(f)

	go func() {
		var event service.Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			args := strings.Split(scanner.Text(), " ")
			for len(args) < 3 {
				args = append(args, "")
			}

			num, err := strconv.Atoi(args[0])
			if err != nil {
				continue
			}

			event.EventType = service.EventType(num)
			event.Key = args[1]
			event.Value = args[2]

			outEvent <- event
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func (t *TransactionLogger) Clear() error {
	return os.Remove(t.path)
}
