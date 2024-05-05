package transactionlogger

import (
	"encoding/base64"
	"fmt"

	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (t *TransactionLogger) WriteSet(key, value string, opts models.SetOptions) {
	readOnly := 1
	if !opts.ReadOnly {
		readOnly = 0
	}
	metadata := fmt.Sprintf("%d;%d", readOnly, opts.Level)
	t.events <- Event{EventType: Set, Name: key, Value: value, Metadata: metadata}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{EventType: Delete, Name: key}
}

func (t *TransactionLogger) WriteSetToObject(name string, key string, val string) {
	t.events <- Event{EventType: SetToObject, Name: name + "." + key, Value: val}
}

func (t *TransactionLogger) WriteCreateObject(name string, info models.ObjectInfo) {
	value := fmt.Sprintf("%d;%d", info.Server, info.Level)
	t.events <- Event{EventType: CreateObject, Name: name, Value: value}
}

func (t *TransactionLogger) WriteDeleteObject(name string) {
	t.events <- Event{EventType: DeleteObject, Name: name}
}

func (t *TransactionLogger) WriteAttach(dst string, src string) {
	t.events <- Event{EventType: Attach, Name: dst, Value: src}
}

func (t *TransactionLogger) WriteDeleteAttr(object string, key string) {
	t.events <- Event{EventType: DeleteAttr, Name: object + "." + key}
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
