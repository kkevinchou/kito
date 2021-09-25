package events

import (
	"encoding/json"

	"github.com/kkevinchou/kito/kito/types"
)

type CreateEntityEvent struct {
	EntityType types.EntityType `json:"entity_type"`
	EntityID   int              `json:"entity_id"`
}

func (e *CreateEntityEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *CreateEntityEvent) Type() EventType {
	return EventTypeCreateEntity
}

func (e *CreateEntityEvent) TypeAsInt() int {
	return int(EventTypeCreateEntity)
}
