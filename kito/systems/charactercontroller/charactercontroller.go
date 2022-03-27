package charactercontroller

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/kito/utils/controllerutils"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayerEntity() entities.Entity
	GetPlayer() *player.Player
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
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()
	singleton := s.world.GetSingleton()

	var players []*player.Player
	if utils.IsClient() {
		players = []*player.Player{s.world.GetPlayer()}
	} else {
		players = playerManager.GetPlayers()
	}

	for _, player := range players {
		entity, err := s.world.GetEntityByID(player.EntityID)
		if err != nil {
			fmt.Printf("error in character controller getting entity %s", err)
			continue
		}

		cameraID := entity.GetComponentContainer().ThirdPersonControllerComponent.CameraID
		camera, err := s.world.GetEntityByID(cameraID)
		if err != nil {
			fmt.Printf("error in character controller getting camera %s", err)
			continue
		}

		controllerutils.UpdateCharacterController(delta, entity, camera, singleton.PlayerInput[player.ID])
	}
}
