package charactercontroller

import (
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/kito/utils/controllerutils"
)

const (
	// a value of 1 means the normal vector of what you're on must be exactly Vec3{0, 1, 0}
	groundedStrictness = 0.85
)

type CharacterControllerResolverSystem struct {
	*base.BaseSystem
	world World
}

func NewCharacterControllerResolverSystem(world World) *CharacterControllerResolverSystem {
	return &CharacterControllerResolverSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CharacterControllerResolverSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	var players []*player.Player
	if utils.IsClient() {
		players = []*player.Player{s.world.GetPlayer()}
	} else {
		players = playerManager.GetPlayers()
	}

	for _, player := range players {
		entity, err := s.world.GetEntityByID(player.EntityID)
		if err != nil {
			continue
		}
		controllerutils.ResolveControllerCollision(entity)

		// cc := entity.GetComponentContainer()
		// capsule := cc.ColliderComponent.CapsuleCollider.Transform(cc.TransformComponent.Position)
		// cc.ColliderComponent.TransformedCapsuleCollider = &capsule

		// cc.ColliderComponent.CollisionInstances

		// if collision.CheckCollisionCapsuleTriangle(capsule, ) {

		// }
	}
}
