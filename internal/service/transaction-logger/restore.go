package transactionlogger

import (
	"fmt"
	"strconv"
	"strings"

	"itisadb/internal/constants"
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
			case 0:
				continue
			case Set:
				split := strings.Split(e.Metadata, constants.MetadataSeparator)

				if len(split) < 2 {
					return fmt.Errorf("%w\n invalid metadata %s, Name: %s", ErrCorruptedConfigFile, e.Metadata, e.Name)
				}

				readOnly := split[0] == "1"

				levelStr := split[1]
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return fmt.Errorf("%w\n invalid level %s, Name: %s", ErrCorruptedConfigFile, levelStr, e.Name)
				}

				if len(split) > 2 {
					encrypt := split[2] == _enctyptedSign
					if encrypt {
						e.Value, err = t.security.Decrypt(e.Value)
						if err != nil {
							return fmt.Errorf("can't decrypt encrypted value %s: %w", e.Name, err)
						}
					}
				}

				r := r.Set(e.Name, e.Value, models.SetOptions{
					ReadOnly: readOnly,
					Level:    models.Level(level),
				})
				if r.IsErr() {
					return fmt.Errorf("can't set %s: %w", e.Name, r.Error())
				}
			case Delete:
				r := r.Delete(e.Name)
				if r.IsErr() {
					return fmt.Errorf("can't delete %s: %w", e.Name, r.Error())
				}
			case SetToObject:
				split := strings.Split(e.Name, constants.ObjectSeparator)
				if len(split) < 2 {
					return fmt.Errorf("%w\n invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}
				key, value := split[len(split)-1], e.Value

				if len(split) > 2 {
					encrypt := split[2] == _enctyptedSign
					if encrypt {
						value, err = t.security.Decrypt(value)
						if err != nil {
							return fmt.Errorf("can't decrypt encrypted value %s: %w", e.Name, err)
						}
					}
				}

				r := r.SetToObject(strings.Join(split[:len(split)-1], constants.ObjectSeparator), key, value, models.SetToObjectOptions{
					ReadOnly: e.Metadata == "1",
				})
				if r.IsErr() {
					return fmt.Errorf("can't set to object %s, v: %s: %w", e.Name, e.Value, r.Error())
				}

			case DeleteAttr:
				r := r.DeleteAttr(e.Name, e.Value)
				if r.IsErr() {
					return fmt.Errorf("can't delete attr %s: %w", e.Name, r.Error())
				}
			case Attach:
				r := r.AttachToObject(e.Name, e.Value)
				if r.IsErr() {
					return fmt.Errorf("can't attach %s, v: %s: %w", e.Name, e.Value, r.Error())
				}
			case DeleteObject:
				rDel := r.DeleteObject(e.Name)
				if rDel.IsErr() {
					return fmt.Errorf("can't delete object %s: %w", e.Name, rDel.Error())
				}
				r.DeleteObjectInfo(e.Name)
				// TODO: case Detach:
			case CreateUser:
				split := strings.Split(e.Metadata, constants.MetadataSeparator)
				if len(split) < 2 {
					return fmt.Errorf("[%w]\n NewUser invalid value %s, Name: %s", ErrCorruptedConfigFile, e.Value, e.Name)
				}

				syncIDStr := split[0]
				activeStr := split[1]
				levelStr := split[2]

				active, err := strconv.ParseBool(activeStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid active value %s, Name: %s", ErrCorruptedConfigFile, e.Metadata, e.Name)
				}

				level, err := strconv.Atoi(levelStr)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid level value %s, Name: %s", ErrCorruptedConfigFile, e.Metadata, e.Name)
				}

				syncID, err := strconv.ParseUint(syncIDStr, 10, 64)
				if err != nil {
					return fmt.Errorf("[%w]\n invalid syncID value %s, Name: %s", ErrCorruptedConfigFile, e.Metadata, e.Name)
				}

				user := models.User{
					Login:    e.Name,
					Password: e.Value,
					Level:    models.Level(level),
					Active:   active,
				}

				user.SetChangeID(syncID)

				rUser := r.NewUser(user)
				if rUser.IsErr() {
					return fmt.Errorf("can't create user %s, v: %s: %w", e.Name, e.Value, rUser.Error())
				}
			case CreateObject:
				split := strings.Split(e.Value, constants.MetadataSeparator)
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

				objOpts := models.ObjectOptions{
					Server: int32(server),
					Level:  models.Level(level),
				}

				rObj := r.CreateObject(e.Name, objOpts)
				if rObj.IsErr() {
					return fmt.Errorf("can't create object %s: %w", e.Name, rObj.Error())
				}

				r.AddObjectInfo(e.Name, models.ObjectInfo(objOpts))
			default:
				return fmt.Errorf("[%w]\n unknown event type %v", ErrCorruptedConfigFile, e)
			}
		}
	}
	return nil
}

func (t *TransactionLogger) Restore(r domains.Restorer) error {
	events, errs := t.readEvents()
	return t.handleEvents(r, events, errs)
}
