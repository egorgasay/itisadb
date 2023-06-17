package transactionlogger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"modernc.org/strutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var PATH = "transaction-logger"

type restorer interface {
	Set(key, value string, uniques bool) error
	Delete(key string) error
	DeleteIndex(name string) error
	CreateIndex(name string) error
	AttachToIndex(dst, src string) error
}

func (t *TransactionLogger) handleIndexes(r restorer, events <-chan Event, errs <-chan error) {
	e, ok := Event{}, true
	var err error

	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.EventType {
			case SetToIndex:
				r.Set(e.Name, e.Value, false)
			case DeleteAttr:
				r.Delete(e.Name)
			case CreateIndex:
				r.CreateIndex(e.Name)
			case Attach:
				r.AttachToIndex(e.Name, e.Value)
			case DeleteIndex:
				r.DeleteIndex(e.Name)
				// TODO: case Detach:
			}
		}
	}
}

func (t *TransactionLogger) handleKV(r restorer, events <-chan Event, errs <-chan error) {
	e, ok := Event{}, true
	var err error

	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.EventType {
			case Set:
				r.Set(e.Name, e.Value, false)
			case Delete:
				r.Delete(e.Name)
			}
		}
	}
}

func (t *TransactionLogger) Restore(r restorer) {
	kvEvents, kvErrs := t.readEvents()
	iEvents, iErrs := t.readEvents()

	go t.handleKV(r, kvEvents, kvErrs)
	go t.handleIndexes(r, iEvents, iErrs)
}

type operation struct {
	sb      strings.Builder
	counter int8
	path    string
}

func newOperation(path string) *operation {
	return &operation{
		path:    path,
		sb:      strings.Builder{},
		counter: 1,
	}
}

func (o *operation) write(data []byte, path string, errs chan<- error, save func(writeTo writer, data string)) {
	o.sb.Write(data)
	o.sb.WriteByte('\n')
	if o.counter == 20 {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		save(f, o.sb.String())
		f.Sync()

		o.sb.Reset()
		o.counter = 0
	}
	o.counter++
}

func (t *TransactionLogger) Run() {
	events := make(chan Event, 60000)
	errorsch := make(chan error, 60000)
	t.errors = errorsch
	t.events = events

	var err error

	n := strconv.Itoa(1)
	t.path = t.path + "/" + n

	op := newOperation(t.path)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var done = make(chan struct{})

	count := 0

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if count < 100_000 {
					continue
				}
				num := 0

				split := strings.Split(t.path, "/")
				if len(split) == 0 {
					continue
				}

				num, err = strconv.Atoi(n)
				if err != nil {
					errorsch <- err
				}
				num++

				t.path = split[0] + "/" + strconv.Itoa(num)
			}
		}
	}()

	go func() {
		defer close(done)
		defer close(errorsch)

		for e := range events {
			data := strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", e.EventType, e.Name, e.Value)))
			op.write(data, t.path, errorsch, t.flash)
			count++
		}
	}()
}

type writer interface {
	WriteString(string) (int, error)
}

// flash grabs 20 events and saves them to the storage.
func (t *TransactionLogger) flash(writeTo writer, data string) {
	t.Lock()
	defer t.Unlock()

	_, err := writeTo.WriteString(data)
	if err != nil {
		log.Println("flash: ", err)
		t.errors <- err
	}
}

func (t *TransactionLogger) Err() <-chan error {
	return t.errors
}

func (t *TransactionLogger) readEventsFrom(r io.Reader, outEvent chan<- Event, outError chan<- error) {
	var event Event
	scanner := bufio.NewScanner(r)
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
		event.Name = args[1]
		event.Value = args[2]

		outEvent <- event
	}

	if err := scanner.Err(); err != nil {
		outError <- fmt.Errorf("transaction log read failure: %w", err)
		return
	}
}

func (t *TransactionLogger) readEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event, 60000)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		d, err := os.ReadDir(t.path)
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}

		for _, f := range d {
			if f.IsDir() {
				continue
			}

			func() {
				file, err := os.Open(t.path + "/" + f.Name())
				if err != nil {
					outError <- fmt.Errorf("transaction log read failure: %w", err)
				}
				defer file.Close()

				t.readEventsFrom(file, outEvent, outError)
			}()

		}

	}()

	return outEvent, outError
}

func (t *TransactionLogger) Stop() error {
	close(t.events)
	return nil
}
