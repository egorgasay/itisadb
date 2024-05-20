package transactionlogger

import (
	"encoding/base64"
	"fmt"

	"itisadb/internal/constants"
	"itisadb/internal/models"

	"go.uber.org/zap"
)

const _enctyptedSign = "E"

func (t *TransactionLogger) WriteSet(key, value string, opts models.SetOptions) {
	readOnly := 1
	if !opts.ReadOnly {
		readOnly = 0
	}

	metadata := fmt.Sprintf("%d%s%d", readOnly, constants.MetadataSeparator, opts.Level)
	if opts.Encrypt {
		encrypted, err := t.security.Encrypt(value)
		if err != nil {
			t.logger.Error("failed to encrypt value", zap.Error(err))
		} else {
			metadata += constants.MetadataSeparator + _enctyptedSign
			value = encrypted
		}
	}

	t.events <- Event{EventType: Set, Name: key, Value: value, Metadata: metadata}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{EventType: Delete, Name: key}
}

func (t *TransactionLogger) WriteSetToObject(name string, key string, val string, opts models.SetToObjectOptions) {
	readOnly := 1
	if !opts.ReadOnly {
		readOnly = 0
	}

	metadata := fmt.Sprintf("%d", readOnly)
	if opts.Encrypt {
		metadata += constants.MetadataSeparator + _enctyptedSign
	}

	t.events <- Event{EventType: SetToObject, Name: name + constants.ObjectSeparator + key, Value: val, Metadata: metadata}
}

func (t *TransactionLogger) WriteCreateObject(name string, info models.ObjectInfo) {
	value := fmt.Sprintf("%d%s%d", info.Server, constants.MetadataSeparator, info.Level)
	t.events <- Event{EventType: CreateObject, Name: name, Value: value}
}

func (t *TransactionLogger) WriteDeleteObject(name string) {
	t.events <- Event{EventType: DeleteObject, Name: name}
}

func (t *TransactionLogger) WriteAttach(dst string, src string) {
	t.events <- Event{EventType: Attach, Name: dst, Value: src}
}

func (t *TransactionLogger) WriteDeleteAttr(object string, key string) {
	t.events <- Event{EventType: DeleteAttr, Name: object + constants.ObjectSeparator + key}
}

var b64 = base64.StdEncoding

func (t *TransactionLogger) WriteNewUser(user models.User) {
	meta := fmt.Sprintf("%d%s%t%s%d", 
		user.GetChangeID(), constants.MetadataSeparator, 
		user.Active, constants.MetadataSeparator, 
		user.Level,
	)

	t.events <- Event{EventType: CreateUser, Name: user.Login, Value: user.Password, Metadata: meta}
}

func (t *TransactionLogger) WriteDeleteUser(login string) {
	t.events <- Event{EventType: DeleteUser, Name: login}
}

//func (t *TransactionLogger) WriteDeleteObjectInfo(name string) {
//	t.events <- Event{EventType: DeleteObjectInfo, Name: name}
//}
