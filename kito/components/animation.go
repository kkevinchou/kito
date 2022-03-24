package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/model"
)

type AnimationComponent struct {
	// stateful data that is manipulated by the Animation System
	ElapsedTime         time.Duration
	Pose                map[int]mgl32.Mat4
	AnimationTransforms map[int]mgl32.Mat4
	CurrentAnimation    string

	// these fields are from the loaded animation and should not be modified
	Animation  *model.Animation
	Animations map[string]*model.Animation
}

func (c *AnimationComponent) GetAnimationComponent() *AnimationComponent {
	return c
}

func (c *AnimationComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AnimationComponent = c
}

func (c *AnimationComponent) PlayAnimation(animationName string) {
	if c.CurrentAnimation != animationName {
		if animation, ok := c.Animations[animationName]; ok {
			c.Animation = animation
			c.CurrentAnimation = animationName
			c.ElapsedTime = 0
		}
	}

}
