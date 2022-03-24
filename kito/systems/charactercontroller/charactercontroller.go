package charactercontroller

import (
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
	GetPlayer() entities.Entity
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
		players = []*player.Player{playerManager.GetPlayer(s.world.GetSingleton().PlayerID)}
	} else {
		players = playerManager.GetPlayers()
	}

	for _, player := range players {
		entity, err := s.world.GetEntityByID(player.ID)
		if err != nil {
			continue
		}

		cameraID := entity.GetComponentContainer().ThirdPersonControllerComponent.CameraID
		camera, err := s.world.GetEntityByID(cameraID)
		if err != nil {
			continue
		}

		controllerutils.UpdateCharacterController(delta, entity, camera, singleton.PlayerInput[player.ID])
	}
}
