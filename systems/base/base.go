package base

import "github.com/kkevinchou/kito/managers/eventbroker"

type BaseSystem struct {
}

func NewBaseSystem() *BaseSystem {
	return &BaseSystem{}
}

func (b *BaseSystem) Observe(eventType eventbroker.EventType, event interface{}) {

}
