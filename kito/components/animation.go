package components

import "github.com/kkevinchou/kito/lib/animation"

type AnimationComponent struct {
	Player *animation.AnimationPlayer
}

func (c *AnimationComponent) GetAnimationComponent() *AnimationComponent {
	return c
}

func (c *AnimationComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AnimationComponent = c
}
