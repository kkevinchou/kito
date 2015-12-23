package physics

import (
	"time"

	"github.com/kkevinchou/ant/math/vector"
)

type Positionable interface {
	Position() vector.Vector
	SetPosition(vector.Vector)
}

type PhysicsComponent struct {
	velocity vector.Vector
	mass     float64
	maxSpeed float64
	entity   Positionable
}

func (c *PhysicsComponent) Init(entity Positionable, maxSpeed, mass float64) {
	c.entity = entity
	c.maxSpeed = maxSpeed
	c.mass = mass
}

func (c *PhysicsComponent) Velocity() vector.Vector {
	return c.velocity
}

func (c *PhysicsComponent) SetVelocity(v vector.Vector) {
	c.velocity = v
}

func (c *PhysicsComponent) Mass() float64 {
	return c.mass
}

func (c *PhysicsComponent) SetMass(mass float64) {
	c.mass = mass
}

func (c *PhysicsComponent) MaxSpeed() float64 {
	return c.maxSpeed
}

func (c *PhysicsComponent) SetMaxSpeed(maxSpeed float64) {
	c.maxSpeed = maxSpeed
}

func (c *PhysicsComponent) Update(delta time.Duration) {
	c.entity.SetPosition(c.entity.Position().Add(c.velocity.Scale(float64(delta.Seconds()))))
}
