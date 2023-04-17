package service

type EventType byte

const (
	Set EventType = iota
	Delete
)

type Event struct {
	EventType EventType
	Key       string
	Value     string
}
