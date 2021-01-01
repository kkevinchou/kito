package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/animation"
)

type AnimationComponent struct {
	ElapsedTime   time.Duration
	Pose          map[int]mgl32.Mat4
	AnimatedModel *animation.AnimatedModel // potentially shared across many entities
	Animation     *animation.Animation

	AnimationTransforms map[int]mgl32.Mat4
}

func (c *AnimationComponent) GetAnimationComponent() *AnimationComponent {
	return c
}

func (c *AnimationComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AnimationComponent = c
}
