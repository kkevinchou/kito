package charactercontroller

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/systems/base"
	"github.com/kkevinchou/kito/systems/common"
	"github.com/kkevinchou/kito/utils"
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
	if utils.IsClient() {
		// return
	}

	d := directory.GetDirectory()
	playerManager := d.PlayerManager()
	singleton := s.world.GetSingleton()

	for _, player := range playerManager.GetPlayers() {
		entity, err := s.world.GetEntityByID(player.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		common.UpdateCharacterController(entity, s.world, singleton.PlayerInput[player.ID])
	}
}
