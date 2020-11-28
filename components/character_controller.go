package components

import (
	"time"

	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type CharacterControllerComponent struct {
	entity        interfaces.Controllable
	controlVector vector.Vector3
}

func NewCharacterControllerComponent(entity interfaces.Controllable) *CharacterControllerComponent {
	component := CharacterControllerComponent{
		entity: entity,
	}
	return &component
}

func (c *CharacterControllerComponent) Update(delta time.Duration) {
	if c.controlVector.IsZero() {
		c.entity.SetVelocity(vector.Zero3())
		return
	}

	forwardVector := c.entity.Backward()
	forwardVector = forwardVector.Scale(-c.controlVector.Z)

	rightVector := c.entity.Right()
	rightVector = rightVector.Scale(-c.controlVector.X)

	velocity := forwardVector.Add(rightVector).Add(vector.Vector3{X: 0, Y: c.controlVector.Y, Z: 0}).Normalize().Scale(c.entity.MaxSpeed())
	c.entity.SetVelocity(velocity)
}

func (c *CharacterControllerComponent) SetVelocityDirection(controlVector vector.Vector3) {
	c.controlVector = controlVector
}
