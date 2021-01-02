package components

import "github.com/go-gl/mathgl/mgl64"

type FollowComponent struct {
	FollowTargetEntityID *int
	InitialOffset        mgl64.Vec3
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
