package ant

import "github.com/kkevinchou/ant/lib/math/vector"

type PositionComponent struct {
	position vector.Vector
}

func (c *PositionComponent) Position() vector.Vector {
	return c.position
}

func (c *PositionComponent) SetPosition(position vector.Vector) {
	c.position = position
}
