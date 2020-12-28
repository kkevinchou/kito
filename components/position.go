package components

import (
	"github.com/go-gl/mathgl/mgl64"
)

type PositionComponent struct {
	position mgl64.Vec3
}

func (c *PositionComponent) Position() mgl64.Vec3 {
	return c.position
}

func (c *PositionComponent) SetPosition(position mgl64.Vec3) {
	c.position = position
}
