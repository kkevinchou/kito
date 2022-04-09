package physics

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/kito/utils/physutils"

	"github.com/kkevinchou/kito/kito/entities"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetPlayerEntity() entities.Entity
	QueryEntity(componentFlags int) []entities.Entity
}

type PhysicsSystem struct {
	*base.BaseSystem
	world World
}

func NewPhysicsSystem(world World) *PhysicsSystem {
	return &PhysicsSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *PhysicsSystem) Update(delta time.Duration) {
	// physics simulation is done on the server and the results are synchronized to the client
	if utils.IsClient() {
		return
	}

	for _, entity := range s.world.QueryEntity(components.ComponentFlagPhysics | components.ComponentFlagTransform) {
		physutils.PhysicsStep(delta, entity)
	}
}
