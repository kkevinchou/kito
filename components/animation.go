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

	// these fields are from the loaded animation and should not be modified
	Animation *model.Animation
}

func (c *AnimationComponent) GetAnimationComponent() *AnimationComponent {
	return c
}

func (c *AnimationComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AnimationComponent = c
}
