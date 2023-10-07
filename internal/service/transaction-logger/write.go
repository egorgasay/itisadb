package transactionlogger

import "itisadb/internal/models"

func (t *TransactionLogger) WriteSet(key, value string) {
	t.events <- Event{EventType: Set, Name: key, Value: value}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{EventType: Delete, Name: key}
}

func (t *TransactionLogger) WriteSetToObject(name string, key string, val string) {
	t.events <- Event{EventType: SetToObject, Name: name + "." + key, Value: val}
}

func (t *TransactionLogger) WriteCreateObject(name string) {
	t.events <- Event{EventType: CreateObject, Name: name}
}

func (t *TransactionLogger) WriteDeleteObject(name string) {
	t.events <- Event{EventType: DeleteObject, Name: name}
}

func (t *TransactionLogger) WriteAttach(dst string, src string) {
	t.events <- Event{EventType: Attach, Name: dst, Value: src}
}

func (t *TransactionLogger) WriteDeleteAttr(name string, key string) {
	t.events <- Event{EventType: DeleteAttr, Name: name + "." + key}
}

func (t *TransactionLogger) WriteCreateUser(user models.User) {
	t.events <- Event{EventType: CreateUser, Name: user.Username, Value: user.Password}
}
