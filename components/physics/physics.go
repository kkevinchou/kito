package physics

import (
	"time"

	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type PhysicsComponent struct {
	velocity vector.Vector3
	mass     float64
	maxSpeed float64
	entity   interfaces.Positionable
	heading  vector.Vector3
}

func (c *PhysicsComponent) Init(entity interfaces.Positionable, maxSpeed, mass float64) {
	c.entity = entity
	c.maxSpeed = maxSpeed
	c.mass = mass
	c.heading = vector.Vector3{X: 0, Y: 0, Z: -1}
}

func (c *PhysicsComponent) Velocity() vector.Vector3 {
	return c.velocity
}

func (c *PhysicsComponent) SetVelocity(v vector.Vector3) {
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
	zero := vector.Vector3{}
	if c.velocity != zero {
		c.heading = c.velocity
	}

	c.entity.SetPosition(c.entity.Position().Add(c.velocity.Scale(float64(delta.Seconds()))))
}

func (c *PhysicsComponent) Heading() vector.Vector3 {
	return c.heading
}
