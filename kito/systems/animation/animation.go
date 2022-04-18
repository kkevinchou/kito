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
		playerEntity := s.world.GetPlayerEntity()
		for _, entity := range s.world.QueryEntity(components.ComponentFlagAnimation) {
			componentContainer := entity.GetComponentContainer()
			animationComponent := componentContainer.AnimationComponent
			player := animationComponent.Player

			tpcComponent := componentContainer.ThirdPersonControllerComponent

			if entity.GetID() == playerEntity.GetID() {
				targetAnimation := "Idle"
				if !libutils.Vec3IsZero(tpcComponent.Velocity) {
					if tpcComponent.Grounded {
						targetAnimation = "Walk"
					} else {
						targetAnimation = "Falling"
					}
				}
				player.PlayAnimation(targetAnimation)
			}
			player.Update(delta)
		}
	} else {
		for _, entity := range s.world.QueryEntity(components.ComponentFlagAnimation) {
			componentContainer := entity.GetComponentContainer()
			animationComponent := componentContainer.AnimationComponent
			player := animationComponent.Player

			tpcComponent := componentContainer.ThirdPersonControllerComponent

			targetAnimation := "Idle"
			if !libutils.Vec3IsZero(tpcComponent.Velocity) {
				if tpcComponent.Grounded {
					targetAnimation = "Walk"
				} else {
					targetAnimation = "Falling"
				}
			}
			player.PlayAnimation(targetAnimation)
			player.Update(delta)
		}
	}
}
