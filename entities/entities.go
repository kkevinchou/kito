package entities

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/settings"
)

var idCounter int = settings.EntityIDStart

type Entity interface {
	GetID() int
	GetName() string
	GetComponentContainer() *components.ComponentContainer
}

type EntityImpl struct {
	ID                 int
	Name               string
	ComponentContainer *components.ComponentContainer
}

func NewEntity(name string, componentContainer *components.ComponentContainer) *EntityImpl {
	e := EntityImpl{
		ID:                 idCounter,
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
