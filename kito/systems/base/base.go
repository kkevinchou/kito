package base

import "github.com/kkevinchou/kito/kito/events"

type BaseSystem struct {
}

func NewBaseSystem() *BaseSystem {
	return &BaseSystem{}
}

func (b *BaseSystem) Observe(event events.Event) {
}
