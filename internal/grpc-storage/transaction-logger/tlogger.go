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
	kvEvents, kvErrs := t.readEvents(t.kvPath)
	iEvents, iErrs := t.readEvents(t.indexesPath)

	go t.handleKV(r, kvEvents, kvErrs)
	go t.handleIndexes(r, iEvents, iErrs)
}

func (t *TransactionLogger) Run() {
	events := make(chan Event, 60000)
	errorsch := make(chan error, 60000)
	t.errors = errorsch
	t.events = events

	var sbIndexes = strings.Builder{}
	var sbKV = strings.Builder{}

	var setToIndexCount = 1
	var setCount = 1

	go func() {
		defer close(events)
		defer close(errorsch)
		for e := range events {
			switch e.EventType {
			case Set:
				sbKV.Write(strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", e.EventType, e.Name, e.Value))))
				sbKV.WriteByte('\n')
				if setCount == 20 {
					f, err := os.OpenFile(t.kvPath+"/1", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
					// TODO: change to dir
					if err != nil {
						errorsch <- err
						continue
					}

					go func() {
						defer f.Close()
						t.flash(f, sbIndexes.String())
					}()

					sbIndexes.Reset()
					setCount = 1
					continue
				}
				setCount++
			case Delete:
				// panic("TODO:")
			case SetToIndex:
				sbIndexes.Write(strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", e.EventType, e.Name, e.Value))))
				sbIndexes.WriteByte('\n')
				if setToIndexCount == 20 {
					f, err := os.OpenFile(t.indexesPath+"/"+e.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						errorsch <- err
						continue
					}

					func() {
						defer f.Close()
						t.flash(f, sbIndexes.String())
					}()

					sbIndexes.Reset()
					setToIndexCount = 1
					continue
				}
				setToIndexCount++
			case DeleteAttr:
				// panic("TODO:")
			case CreateIndex:
				f, err := os.Create(t.indexesPath + "/" + e.Name)
				if err != nil {
					errorsch <- err
					continue
				}
				f.Close()

			case Attach:
				// panic("TODO:")
			case DeleteIndex:
				if err := os.Remove(t.indexesPath + "/" + e.Name); err != nil {
					errorsch <- err
					continue
				}
			}
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

func (t *TransactionLogger) readEvents(dir string) (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		d, err := os.ReadDir(dir)
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}

		for _, f := range d {
			if f.IsDir() {
				continue
			}

			func() {
				file, err := os.Open(dir + "/" + f.Name())
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

func (t *TransactionLogger) Clear() error {
	// err := os.RemoveAll(t.path)
	return nil
}
