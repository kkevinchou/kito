package components

import (
	"time"

	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type CharacterControllerComponent struct {
	entity        interfaces.Controllable
	controlVector vector.Vector3
	zoomValue     int
}

func NewCharacterControllerComponent(entity interfaces.Controllable) *CharacterControllerComponent {
	component := CharacterControllerComponent{
		entity: entity,
	}
	return &component
}

func (c *CharacterControllerComponent) Update(delta time.Duration) {
	if c.controlVector.IsZero() && c.zoomValue == 0 {
		c.entity.SetVelocity(vector.Zero3())
		return
	}

	forwardVector := c.entity.Forward()
	zoomVector := forwardVector.Scale(float64(-c.zoomValue))

	forwardVector = forwardVector.Scale(c.controlVector.Z)
	forwardVector.Y = 0

	rightVector := c.entity.Right()
	rightVector = rightVector.Scale(-c.controlVector.X)

	velocity := zoomVector.Add(forwardVector).Add(rightVector).Add(vector.Vector3{X: 0, Y: c.controlVector.Y, Z: 0}).Normalize().Scale(c.entity.MaxSpeed())
	c.entity.SetVelocity(velocity)
}

func (c *CharacterControllerComponent) SetControlDirection(controlVector vector.Vector3, zoom int) {
	c.controlVector = controlVector
	c.zoomValue = zoom
}
