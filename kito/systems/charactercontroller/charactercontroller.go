package charactercontroller

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/netsync"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) entities.Entity
	GetPlayerEntity() entities.Entity
	GetPlayer() *player.Player
}

type CharacterControllerSystem struct {
	*base.BaseSystem
	world World
}

func NewCharacterControllerSystem(world World) *CharacterControllerSystem {
	return &CharacterControllerSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
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
		entity := s.world.GetEntityByID(player.EntityID)
		if entity == nil {
			fmt.Printf("character controller could not find player entity with id %d\n", player.EntityID)
			continue
		}

		cameraID := entity.GetComponentContainer().ThirdPersonControllerComponent.CameraID
		camera := s.world.GetEntityByID(cameraID)
		if camera == nil {
			fmt.Printf("character controller could not find camera with entity id %d\n", cameraID)
			continue
		}

		netsync.UpdateCharacterController(delta, entity, camera, singleton.PlayerInput[player.ID])
	}
}

func (s *CharacterControllerSystem) Name() string {
	return "CharacterControllerSystem"
}
