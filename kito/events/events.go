package events

type EventType string

var EventTypeUnregisterEntity EventType = "UNREGISTER"

type Event interface {
	Type() EventType
}
