package transactionlogger

import (
	"fmt"
	"strings"
)

type restorer interface {
	Set(key, value string, uniques bool) error
	Delete(key string) error
	SetToObject(name, key, value string, uniques bool) error
	DeleteObject(name string) error
	CreateObject(name string) error
	AttachToObject(dst, src string) error
}

var ErrCorruptedConfigFile = fmt.Errorf("corrupted config file")

func (t *TransactionLogger) handleEvents(r restorer, events <-chan Event, errs <-chan error) error {
	e, ok := Event{}, true
	var err error

	for ok && err == nil {
		select {
		case err, ok = <-errs:
			if ok && err != nil {
				return err
			}
			ok = true
		case e, ok = <-events:
			switch e.EventType {
			case Set:
				r.Set(e.Name, e.Value, false)
			case Delete:
				r.Delete(e.Name)
			case SetToObject:
				split := strings.Split(e.Value, ".")
				if len(split) != 2 {
					return fmt.Errorf("%w\n invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}
				key, value := split[0], split[1]
				r.SetToObject(e.Name, key, value, false)
			case DeleteAttr:
				r.Delete(e.Name)
			case CreateObject:
				r.CreateObject(e.Name)
			case Attach:
				r.AttachToObject(e.Name, e.Value)
			case DeleteObject:
				r.DeleteObject(e.Name)
				// TODO: case Detach:
			}
		}
	}
	return nil
}

func (t *TransactionLogger) Restore(r restorer) error {
	events, errs := t.readEvents()
	return t.handleEvents(r, events, errs)
}
