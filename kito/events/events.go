package events

type EventType string

var EventTypeUnregisterEntity EventType = "UNREGISTER"
var EventTypeConsoleEnabled EventType = "CONSOLE_ENABLED"
var EventTypeRPC EventType = "RPC"

type Event interface {
	Type() EventType
}

type UnregisterEntityEvent struct {
	EntityID           int `json:"entity_id"`
	GlobalCommandFrame int `json:"global_command_frame"`
}

func (e *UnregisterEntityEvent) Type() EventType {
	return EventTypeUnregisterEntity
}

type ConsoleEnabledEvent struct {
}

func (e *ConsoleEnabledEvent) Type() EventType {
	return EventTypeConsoleEnabled
}

type RPCEvent struct {
	Command string
}

func (e *RPCEvent) Type() EventType {
	return EventTypeRPC
}
