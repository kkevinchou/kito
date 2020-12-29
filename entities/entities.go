package entities

import "github.com/kkevinchou/kito/components"

type Entity interface {
	GetComponentContainer() *components.ComponentContainer
}

type EntityImpl struct {
	ComponentContainer *components.ComponentContainer
}

func (e *EntityImpl) GetComponentContainer() *components.ComponentContainer {
	return e.ComponentContainer
}
