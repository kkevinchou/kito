package charactercontroller

import (
	"time"

	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/kkevinchou/kito/systems/common"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type CharacterControllerSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCharacterControllerSystem(world World) *CharacterControllerSystem {
	return &CharacterControllerSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CharacterControllerSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil && componentContainer.ThirdPersonControllerComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CharacterControllerSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	for _, entity := range s.entities {
		common.UpdateCharacterController(entity, s.world, singleton.Input)
		// fmt.Println(entity.GetComponentContainer().TransformComponent.Position)
	}
}
