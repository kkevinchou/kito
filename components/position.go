package components

import "github.com/kkevinchou/ant/lib/math/vector"

type PositionComponent struct {
	position vector.Vector3
}

func (c *PositionComponent) Position() vector.Vector3 {
	return c.position
}

func (c *PositionComponent) SetPosition(position vector.Vector3) {
	c.position = position
}
