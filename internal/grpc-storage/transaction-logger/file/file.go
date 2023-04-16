package file

import (
	"bufio"
	"fmt"
	"grpc-storage/internal/grpc-storage/transaction-logger/service"
	"log"
	"modernc.org/strutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

type TransactionLogger struct {
	path string
	file *os.File

	events chan service.Event
	errors chan error

	sync.RWMutex
}

func NewLogger(path string) (*TransactionLogger, error) {
	path = strings.TrimRight(path, "/") + "/transactionLogger"

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

	t.events = events

	var sb = strings.Builder{}
	var count = 1

	go func() {
		for e := range events {
			sb.Write(strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", e.EventType, e.Key, e.Value))))
			sb.WriteByte('\n')
			if count == 20 {
				go t.flash(sb.String())
				sb.Reset()
				count = 1
				continue
			}
			count++
		}
	}()
}

// flash grabs 20 events and saves them to the db.
func (t *TransactionLogger) flash(data string) {
	t.Lock()
	defer t.Unlock()
	errorsch := make(chan error, 20)
	t.errors = errorsch

	_, err := t.file.WriteString(data)
	if err != nil {
		log.Println("flash: ", err)
		errorsch <- err
	}
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
			action := scanner.Text()
			decode, err := strutil.Base64Decode([]byte(action))
			if err != nil {
				return
			}

			args := strings.Split(string(decode), " ")
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
