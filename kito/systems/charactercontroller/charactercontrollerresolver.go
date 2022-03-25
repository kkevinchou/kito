package charactercontroller

import (
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/kito/utils/controllerutils"

	"github.com/kkevinchou/kito/kito/entities"
)

const (
	// a value of 1 means the normal vector of what you're on must be exactly Vec3{0, 1, 0}
	groundedStrictness = 0.85
)

type CharacterControllerResolverSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCharacterControllerResolverSystem(world World) *CharacterControllerResolverSystem {
	return &CharacterControllerResolverSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		entities:   []entities.Entity{},
	}
}

func (s *CharacterControllerResolverSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.ThirdPersonControllerComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}

}

func (s *CharacterControllerResolverSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

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
		controllerutils.ResolveControllerCollision(entity)
	}
}
