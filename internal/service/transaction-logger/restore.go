package transactionlogger

import (
	"fmt"
	"strconv"
	"strings"

	"itisadb/internal/domains"
	"itisadb/internal/models"
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
				split := strings.Split(e.Metadata, ";")

				if len(split) < 2 {
					return fmt.Errorf("%w\n invalid metadata %s, Name: %s", ErrCorruptedConfigFile, e.Metadata, e.Name)
				}

				readOnly := split[0] == "1"

				levelStr := split[1]
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return fmt.Errorf("%w\n invalid level %s, Name: %s", ErrCorruptedConfigFile, levelStr, e.Name)
				}

				r.Set(e.Name, e.Value, models.SetOptions{
					ReadOnly: readOnly,
					Level:    models.Level(level),
				})
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
				split := strings.Split(e.Metadata, ";")
				if len(split) < 2 {
					return fmt.Errorf("[%w]\n CreateUser invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				activeStr := split[0]
				levelStr := split[1]

				active, err := strconv.ParseBool(activeStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid active value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid level value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				r.CreateUser(models.User{
					Login:    e.Name,
					Password: e.Value,
					Level:    models.Level(level),
					Active:   active,
				})
			case AddObjectInfo:
				split := strings.Split(e.Value, ";")
				if len(split) < 2 {
					return fmt.Errorf("[%w]\n AddObjectInfo invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				serverStr := split[0]
				levelStr := split[1]

				server, err := strconv.Atoi(serverStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid server value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid level value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				r.AddObjectInfo(e.Name, models.ObjectInfo{
					Server: int32(server),
					Level:  models.Level(level),
				})
			case DeleteObjectInfo:
				r.DeleteObjectInfo(e.Name)
			}
		}
	}
	return nil
}

func (t *TransactionLogger) Restore(r domains.Restorer) error {
	events, errs := t.readEvents()
	return t.handleEvents(r, events, errs)
}
