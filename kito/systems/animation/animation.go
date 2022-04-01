package animation

import (
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/libutils"
)

type World any

type AnimationSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewAnimationSystem(world World) *AnimationSystem {
	return &AnimationSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		entities:   []entities.Entity{},
	}
}

func (s *AnimationSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.AnimationComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *AnimationSystem) Update(delta time.Duration) {
	for _, entity := range s.entities {
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
