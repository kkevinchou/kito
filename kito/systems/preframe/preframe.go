package preframe

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
	GetSingleton() *singleton.Singleton
	QueryEntity(componentFlags int) []entities.Entity
}

type PreFrameSystem struct {
	*base.BaseSystem

	world World
}

func NewPreFrameSystem(world World) *PreFrameSystem {
	return &PreFrameSystem{
		world: world,
	}
}

func (s *PreFrameSystem) Update(delta time.Duration) {
	if utils.IsServer() {
		singleton := s.world.GetSingleton()
		for i, _ := range singleton.PlayerInput {
			singleton.PlayerInput[i] = input.Input{}
		}
	}

	// reset notepad
	for _, entity := range s.world.QueryEntity(components.ComponentFlagNotepad) {
		entity.GetComponentContainer().NotepadComponent.LastAction = components.ActionNone
	}

}

func (s *PreFrameSystem) Name() string {
	return "PreFrameSystem"
}
