package physics

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/math"
	"github.com/kkevinchou/kito/types"
)

const (
	fullDecayThreshold = float64(0.05)
)

type PhysicsComponent struct {
	velocity mgl64.Vec3
	mass     float64
	maxSpeed float64
	entity   types.Positionable
	heading  mgl64.Vec3

	// impulses have a name that can be reset or overwritten
	impulses map[string]*types.Impulse
}

func (c *PhysicsComponent) Init(entity types.Positionable, maxSpeed, mass float64) {
	c.entity = entity
	c.maxSpeed = maxSpeed
	c.mass = mass
	c.heading = mgl64.Vec3{0, 0, -1}
	c.impulses = map[string]*types.Impulse{}
}

func (c *PhysicsComponent) Velocity() mgl64.Vec3 {
	return c.velocity
}

func (c *PhysicsComponent) SetVelocity(v mgl64.Vec3) {
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

func (c *PhysicsComponent) ApplyImpulse(name string, impulse *types.Impulse) {
	c.impulses[name] = impulse
}

func (c *PhysicsComponent) Update(delta time.Duration) {
	if !math.Vec3IsZero(c.velocity) {
		c.heading = c.velocity
	}

	var totalImpulse mgl64.Vec3
	for name := range c.impulses {
		impulse := c.impulses[name]
		impulse.ElapsedTime = impulse.ElapsedTime + delta
		decayRatio := 1.0 - (impulse.ElapsedTime.Seconds() * impulse.DecayRate)
		if decayRatio < 0 {
			decayRatio = 0
		}

		if decayRatio < fullDecayThreshold {
			delete(c.impulses, name)
		} else {
			realImpulse := impulse.Vector.Mul(decayRatio)
			totalImpulse = totalImpulse.Add(realImpulse)
		}
	}

	velocity := c.velocity.Add(totalImpulse)
	newPos := c.entity.Position().Add(velocity.Mul(delta.Seconds()))
	c.entity.SetPosition(newPos)
}

func (c *PhysicsComponent) Heading() mgl64.Vec3 {
	return c.heading
}
