package animation

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/libutils"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
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

func (s *AnimationSystem) RegisterEntity(entity entities.Entity) {
}

func (s *AnimationSystem) Update(delta time.Duration) {
	for _, entity := range s.world.QueryEntity(components.ComponentFlagAnimation) {
		componentContainer := entity.GetComponentContainer()
		animationComponent := componentContainer.AnimationComponent
		tpcComponent := componentContainer.ThirdPersonControllerComponent

		targetAnimation := "Idle"
		if !libutils.Vec3IsZero(tpcComponent.Velocity) {
			targetAnimation = "Walk"
		}

		player := animationComponent.Player
		player.PlayAnimation(targetAnimation)
		player.Update(delta)
	}
}
