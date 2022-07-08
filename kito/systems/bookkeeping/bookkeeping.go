package bookkeeping

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	QueryEntity(componentFlags int) []entities.Entity
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

func (s *BookKeepingSystem) Update(delta time.Duration) {
	if utils.IsClient() {
		for _, entity := range s.world.QueryEntity(components.ComponentFlagNotepad) {
			entity.GetComponentContainer().NotepadComponent.GetNotepadComponent().LastAction = components.ActionNone
		}
	} else {
		singleton := s.world.GetSingleton()
		for i, _ := range singleton.PlayerInput {
			singleton.PlayerInput[i] = input.Input{}
		}
	}
}
