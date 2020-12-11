package animation

import "time"

type World interface{}

type AnimationSystem struct {
	world World
}

func NewAnimationSystem(world World) *AnimationSystem {
	return &AnimationSystem{
		world: world,
	}
}

func (s *AnimationSystem) Update(delta time.Duration) {
}
