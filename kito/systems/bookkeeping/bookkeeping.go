package bookkeeping

import (
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
}

type BookKeepingSystem struct {
	*base.BaseSystem

	world World
}

func NewBookKeepingSystem(world World) *BookKeepingSystem {
	return &BookKeepingSystem{
		world: world,
	}
}

func (s *BookKeepingSystem) RegisterEntity(entity entities.Entity) {
}

func (s *BookKeepingSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	for i, _ := range singleton.PlayerInput {
		singleton.PlayerInput[i] = input.Input{}
	}
}
