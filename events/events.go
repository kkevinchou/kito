package events

type EventType int

type Event interface {
	Type() EventType
	Serialize() ([]byte, error)
}

const (
	EventTypeCreateEntity EventType = iota
)
