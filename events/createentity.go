package events

import "encoding/json"

type EntityType int

const (
	EntityTypeBob EntityType = iota
)

type CreateEntityEvent struct {
	EntityType EntityType `json:"entity_type"`
	EntityID   int        `json:"entity_id"`
}

func (e *CreateEntityEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *CreateEntityEvent) Type() EventType {
	return EventTypeCreateEntity
}
