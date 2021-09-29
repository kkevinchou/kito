package events

type EventType int

type Event interface {
	Type() EventType
	TypeAsInt() int
	Serialize() ([]byte, error)
}

const (
	EventTypeCreateEntity EventType = iota
)