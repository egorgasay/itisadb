package transactionlogger

func (t *TransactionLogger) WriteSet(key, value string) {
	t.events <- Event{EventType: Set, Name: key, Value: value}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{EventType: Delete, Name: key}
}

func (t *TransactionLogger) WriteSetToIndex(name string, key string, val string) {
	t.events <- Event{EventType: SetToIndex, Name: name + "." + key, Value: val}
}

func (t *TransactionLogger) WriteCreateIndex(name string) {
	t.events <- Event{EventType: CreateIndex, Name: name}
}

func (t *TransactionLogger) WriteDeleteIndex(name string) {
	t.events <- Event{EventType: DeleteIndex, Name: name}
}

func (t *TransactionLogger) WriteAttach(dst string, src string) {
	t.events <- Event{EventType: Attach, Name: dst, Value: src}
}

func (t *TransactionLogger) WriteDeleteAttr(name string, key string) {
	t.events <- Event{EventType: DeleteAttr, Name: name + "." + key}
}
