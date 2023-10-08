package transactionlogger

import (
	"fmt"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"strings"
)

var ErrCorruptedConfigFile = fmt.Errorf("corrupted config file")

func (t *TransactionLogger) handleEvents(r domains.Restorer, events <-chan Event, errs <-chan error) error {
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
				r.Set(e.Name, e.Value, models.SetOptions{})
			case Delete:
				r.Delete(e.Name)
			case SetToObject:
				split := strings.Split(e.Name, ".")
				if len(split) < 2 {
					return fmt.Errorf("%w\n invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}
				key, value := split[len(split)-1], e.Value
				r.SetToObject(strings.Join(split[:len(split)-1], "."), key, value, models.SetToObjectOptions{})
			case DeleteAttr:
				r.Delete(e.Name)
			case CreateObject:
				r.CreateObject(e.Name, models.ObjectOptions{})
			case Attach:
				r.AttachToObject(e.Name, e.Value)
			case DeleteObject:
				r.DeleteObject(e.Name)
				// TODO: case Detach:
			case CreateUser:
				r.CreateUser(models.User{
					Username: e.Name,
					Password: e.Value,
				})

			}
		}
	}
	return nil
}

func (t *TransactionLogger) Restore(r domains.Restorer) error {
	events, errs := t.readEvents()
	return t.handleEvents(r, events, errs)
}
