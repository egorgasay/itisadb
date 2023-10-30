package transactionlogger

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"modernc.org/strutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var PATH = "transaction-logger"

const MaxBufferSize = 20
const MaxCOL = 100_000

type limitedBuffer struct {
	sb      strings.Builder
	counter int16
}

func newLimitedBuffer() *limitedBuffer {
	return &limitedBuffer{
		sb:      strings.Builder{},
		counter: 1,
	}
}

func (t *TransactionLogger) Run() {
	events := make(chan Event, 60000)
	errorsch := make(chan error, 60000)

	t.errors = errorsch
	t.events = events

	op := newLimitedBuffer()
	done := make(chan struct{})

	go t.countWatcher(done)

	go func() {
		defer close(done)
		defer close(errorsch)

		for e := range events {
			var data []byte
			// TODO: Check this
			data = []byte(
				fmt.Sprintf(
					"%d %s %s %s\n",
					e.EventType,
					b64.EncodeToString([]byte(e.Name)),
					b64.EncodeToString([]byte(e.Value)),
					b64.EncodeToString([]byte(e.Metadata)),
				),
			)

			op.sb.Write(data)

			op.counter++
			t.currentCOL++

			if op.counter == MaxBufferSize {
				t.RLock()
				_, err := t.file.WriteString(op.sb.String())
				//t.file.Sync() TODO: ???
				t.RUnlock()
				if err != nil {
					t.logger.Error("flash error", zap.Error(err))
					t.errors <- err
				}
				op.sb.Reset()
				op.counter = 0
			}
		}
	}()
}

func (t *TransactionLogger) countWatcher(done chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if t.currentCOL < MaxCOL {
				continue
			}
			t.currentName++

			t.pathToFile = fmt.Sprintf("%s/%d", PATH, t.currentName)
			f, err := os.OpenFile(t.pathToFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.errors <- err
			}

			t.Lock()
			t.currentCOL = 0
			t.file.Close()
			t.file = f
			t.Unlock()
		}
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
		event.Value = strings.TrimSpace(strings.Join(args[2:], " "))

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

		d, err := os.ReadDir(t.pathToDir)
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}

		i := 1
		for _, f := range d {
			if f.IsDir() {
				continue
			}

			func() {
				file, err := os.Open(t.pathToDir + "/" + fmt.Sprintf("%d", i))
				if err != nil {
					outError <- fmt.Errorf("transaction log read failure: %w", err)
				}
				defer file.Close()

				t.readEventsFrom(file, outEvent, outError)
				i++
			}()

		}

	}()

	return outEvent, outError
}

func (t *TransactionLogger) Stop() error {
	close(t.events)
	return t.file.Close()
}
