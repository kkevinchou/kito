package components

import "github.com/go-gl/mathgl/mgl64"

type FollowComponent struct {
	FollowTargetEntityID *int

	InitialOffset      mgl64.Vec3
	CurrentOffsetDelta mgl64.Vec3

	Rotation mgl64.Quat
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
