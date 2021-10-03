package components

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/types"
)

type PhysicsComponent struct {
	Velocity mgl64.Vec3

	// impulses have a name that can be reset or overwritten
	Impulses map[string]types.Impulse
}

func (c *PhysicsComponent) ApplyImpulse(name string, impulse types.Impulse) {
	c.Impulses[name] = impulse
}

func (c *PhysicsComponent) AddToComponentContainer(container *ComponentContainer) {
	container.PhysicsComponent = c
}
