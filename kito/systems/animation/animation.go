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
	var entities []entities.Entity
	if utils.IsClient() {
		playerEntity := s.world.GetPlayerEntity()
		componentContainer := playerEntity.GetComponentContainer()
		if animationComponent := componentContainer.AnimationComponent; animationComponent != nil {
			entities = append(entities, playerEntity)
		}
	} else {
		entities = s.world.QueryEntity(components.ComponentFlagAnimation)
	}

	playAnimationsForEntities(delta, entities)
}

func playAnimationsForEntities(delta time.Duration, entities []entities.Entity) {
	for _, entity := range entities {
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
