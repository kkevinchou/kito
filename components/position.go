package components

import (
	"github.com/go-gl/mathgl/mgl64"
)

type PositionComponent struct {
	Position mgl64.Vec3
	View     mgl64.Vec3
}

func (c *PositionComponent) AddToComponentContainer(container *ComponentContainer) {
	container.PositionComponent = c
}
