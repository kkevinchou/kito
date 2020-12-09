package physics

import (
	"time"

	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
)

const (
	// the decay fraction per second for an impulse
	impulseDecay       = float64(3)
	fullDecayThreshold = float64(0.05)
)

type Impulse struct {
	v           vector.Vector3
	elpasedTime time.Duration
}

type PhysicsComponent struct {
	velocity vector.Vector3
	mass     float64
	maxSpeed float64
	entity   interfaces.Positionable
	heading  vector.Vector3

	// impulses have a name that can be reset or overwritten
	impulses map[string]*Impulse
}

func (c *PhysicsComponent) Init(entity interfaces.Positionable, maxSpeed, mass float64) {
	c.entity = entity
	c.maxSpeed = maxSpeed
	c.mass = mass
	c.heading = vector.Vector3{X: 0, Y: 0, Z: -1}
	c.impulses = map[string]*Impulse{}
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

func (c *PhysicsComponent) ApplyImpulse(name string, impulse vector.Vector3) {
	c.impulses[name] = &Impulse{v: impulse}
}

func (c *PhysicsComponent) Update(delta time.Duration) {
	if !c.velocity.IsZero() {
		c.heading = c.velocity
	}

	var totalImpulse vector.Vector3
	for name := range c.impulses {
		c.impulses[name].elpasedTime = c.impulses[name].elpasedTime + delta
		decayRatio := 1.0 - (c.impulses[name].elpasedTime.Seconds() * impulseDecay)
		if decayRatio < 0 {
			decayRatio = 0
		}

		if decayRatio < fullDecayThreshold {
			delete(c.impulses, name)
		} else {
			realImpulse := c.impulses[name].v.Scale(decayRatio)
			totalImpulse = totalImpulse.Add(realImpulse)
		}
	}

	velocity := c.velocity.Add(totalImpulse)
	newPos := c.entity.Position().Add(velocity.Scale(float64(delta.Seconds())))
	c.entity.SetPosition(newPos)
}

func (c *PhysicsComponent) Heading() vector.Vector3 {
	return c.heading
}
