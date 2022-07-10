package events

type EventType int

type Event interface {
	Type() EventType
	TypeAsInt() int
	Serialize() ([]byte, error)
}

const (
	EventTypeUnregisterEntity EventType = iota
)
