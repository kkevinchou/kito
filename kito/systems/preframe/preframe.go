package preframe

import (
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
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

}

func (s *PreFrameSystem) Name() string {
	return "PreFrameSystem"
}
