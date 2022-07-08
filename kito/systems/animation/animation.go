package animation

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/libutils"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	GetPlayerEntity() entities.Entity
	GetEntityByID(id int) entities.Entity
}

type AnimationSystem struct {
	*base.BaseSystem
	world World
}

func NewAnimationSystem(world World) *AnimationSystem {
	return &AnimationSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *AnimationSystem) Update(delta time.Duration) {
	if utils.IsClient() {
		// play animations for the player
		playerEntity := s.world.GetPlayerEntity()
		findAndPlayAnimation(delta, playerEntity)
		playerEntity.GetComponentContainer().AnimationComponent.Player.Update(delta)

		// update the animation player for all other entities, relying on animation state
		// synchronization from the server
		for _, entity := range s.world.QueryEntity(components.ComponentFlagAnimation) {
			if entity.GetID() == playerEntity.GetID() {
				continue
			}
			entity.GetComponentContainer().AnimationComponent.Player.Update(delta)
		}
	} else {
		for _, entity := range s.world.QueryEntity(components.ComponentFlagAnimation) {
			findAndPlayAnimation(delta, entity)
			entity.GetComponentContainer().AnimationComponent.Player.Update(delta)
		}
	}
}

// findAndPlayAnimation takes an entity and finds the appropriate animation to play based on its state, then plays it
func findAndPlayAnimation(delta time.Duration, entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()
	animationComponent := componentContainer.AnimationComponent
	player := animationComponent.Player

	tpcComponent := componentContainer.ThirdPersonControllerComponent

	var targetAnimation string
	if !libutils.Vec3IsZero(tpcComponent.Velocity) {
		if tpcComponent.Grounded {
			targetAnimation = "Walk"
		} else {
			targetAnimation = "Falling"
		}
		player.PlayAnimation(targetAnimation)
	} else {
		targetAnimation = "Idle"
		notepad := componentContainer.NotepadComponent
		if notepad.LastAction == components.ActionCast {
			targetAnimation = "Cast1"
			player.PlayOnce(targetAnimation, "Idle")
		} else {
			player.PlayAnimation(targetAnimation)
		}
	}
}
