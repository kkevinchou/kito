package components

import (
	"github.com/go-gl/mathgl/mgl64"
)

type TransformComponent struct {
	Position      mgl64.Vec3
	Orientation   mgl64.Quat
	UpVector      mgl64.Vec3
	ForwardVector mgl64.Vec3
}

func (c *TransformComponent) AddToComponentContainer(container *ComponentContainer) {
	container.TransformComponent = c
}
