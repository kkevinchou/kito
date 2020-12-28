package animation

import (
	"time"

	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

type World interface{}

type AnimationSystem struct {
	world    World
	animator *animation.Animator
}

func NewAnimationSystem(world World) *AnimationSystem {
	parsedCollada, err := collada.ParseCollada("_assets/collada/cube2.dae")
	if err != nil {
		panic(err)
	}
	animatedModel := animation.NewAnimatedModel(parsedCollada, 50, 3)

	return &AnimationSystem{
		world:    world,
		animator: animation.NewAnimator(animatedModel, parsedCollada.Animation),
	}
}

func (s *AnimationSystem) Update(delta time.Duration) {
	s.animator.Update(delta)
}
