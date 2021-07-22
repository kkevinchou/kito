package entities

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/types"
)

var idCounter int = settings.EntityIDStart

type Entity interface {
	GetID() int
	Type() types.EntityType
	GetName() string
	GetComponentContainer() *components.ComponentContainer
}

type EntityImpl struct {
	ID                 int
	entityType         types.EntityType
	Name               string
	ComponentContainer *components.ComponentContainer
}

func NewEntity(name string, entityType types.EntityType, componentContainer *components.ComponentContainer) *EntityImpl {
	e := EntityImpl{
		ID:                 idCounter,
		entityType:         entityType,
		Name:               name,
		ComponentContainer: componentContainer,
	}
	idCounter++
	return &e
}

func (e *EntityImpl) GetComponentContainer() *components.ComponentContainer {
	return e.ComponentContainer
}

func (e *EntityImpl) GetName() string {
	return e.Name
}

func (e *EntityImpl) GetID() int {
	return e.ID
}

func (e *EntityImpl) Type() types.EntityType {
	return e.entityType
}
