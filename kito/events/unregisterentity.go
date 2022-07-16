package events

type UnregisterEntityEvent struct {
	EntityID           int `json:"entity_id"`
	GlobalCommandFrame int `json:"global_command_frame"`
}

func (e *UnregisterEntityEvent) Type() EventType {
	return EventTypeUnregisterEntity
}
