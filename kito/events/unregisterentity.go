package events

import (
	"encoding/json"
)

type UnregisterEntityEvent struct {
	EntityID int `json:"entity_id"`
}

func (e *UnregisterEntityEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UnregisterEntityEvent) Type() EventType {
	return EventTypeUnregisterEntity
}

func (e *UnregisterEntityEvent) TypeAsInt() int {
	return int(EventTypeUnregisterEntity)
}

func DeserializeUnregisterEntityEvent(bytes []byte) UnregisterEntityEvent {
	event := UnregisterEntityEvent{}
	err := json.Unmarshal(bytes, &event)
	if err != nil {
		panic(err)
	}
	return event
}
