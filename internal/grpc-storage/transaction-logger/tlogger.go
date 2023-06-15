package transactionlogger

import (
	"bufio"
	"fmt"
	"log"
	"modernc.org/strutil"
	"os"
	"strconv"
	"strings"
)

func NewLogger(path string) (*TransactionLogger, error) {
	path = strings.TrimRight(path, "/") + "/transactionLogger"

	return &TransactionLogger{
		path: path,
	}, nil
}

func (t *TransactionLogger) WriteSet(key, value string) {
	t.events <- Event{EventType: Set, Key: key, Value: value}
}

type restorer interface {
	Set(string, string, bool)
	Delete(string)
}

func (t *TransactionLogger) Restore(r restorer) {
	var err error

	events, errs := t.readEvents()
	e, ok := Event{}, true

	t.run()

	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.EventType {
			case Set:
				r.Set(e.Key, e.Value, false)
			case Delete:
				r.Delete(e.Key)
			}
		}
	}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{EventType: Delete, Key: key}
}

func (t *TransactionLogger) run() {
	events := make(chan Event, 60000)
	errorsch := make(chan error, 60000)
	t.errors = errorsch
	t.events = events

	var sb = strings.Builder{}
	var count = 1

	go func() {
		defer close(events)
		defer close(errorsch)
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

	_, err := t.file.WriteString(data)
	if err != nil {
		log.Println("flash: ", err)
		t.errors <- err
	}
}

func (t *TransactionLogger) Err() <-chan error {
	return t.errors
}

func (t *TransactionLogger) readEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)
	f, err := os.OpenFile(t.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	t.file = f
	if err != nil {
		outError <- err
		return outEvent, outError
	}
	scanner := bufio.NewScanner(f)

	go func() {
		var event Event
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

			event.EventType = EventType(num)
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
