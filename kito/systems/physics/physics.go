package physics

import (
	"time"

	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/kito/utils/physutils"

	"github.com/kkevinchou/kito/kito/entities"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetPlayer() entities.Entity
}

type PhysicsSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewPhysicsSystem(world World) *PhysicsSystem {
	return &PhysicsSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		entities:   []entities.Entity{},
	}
}

func (s *PhysicsSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *PhysicsSystem) Update(delta time.Duration) {
	if utils.IsClient() {
		player := s.world.GetPlayer()
		if player != nil {
			physutils.PhysicsStep(delta, s.world.GetPlayer())
		}
	} else {
		for _, entity := range s.entities {
			physutils.PhysicsStep(delta, entity)
		}
	}
}
