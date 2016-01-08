package physics

import (
	"time"

	"github.com/kkevinchou/ant/lib/math/vector"
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
	heading  vector.Vector
}

func (c *PhysicsComponent) Init(entity Positionable, maxSpeed, mass float64) {
	c.entity = entity
	c.maxSpeed = maxSpeed
	c.mass = mass
	c.heading = vector.Vector{X: 1, Y: 0}
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
	if c.velocity != vector.Zero() {
		c.heading = c.velocity
	}

	c.entity.SetPosition(c.entity.Position().Add(c.velocity.Scale(float64(delta.Seconds()))))
}

func (c *PhysicsComponent) Heading() vector.Vector {
	return c.heading
}
