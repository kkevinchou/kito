package physics

import (
	// "fmt"

	"time"

	"github.com/kkevinchou/ant/math/vector"
)

type PhysicsI interface {
	GetPhysicsComponent() *PhysicsComponent
	Update(delta time.Duration)
}

type PhysicsComponent struct {
	Position vector.Vector
	Velocity vector.Vector
	Mass     float64
	MaxSpeed float64
}

func (component *PhysicsComponent) Update(delta time.Duration) {
	component.Position = component.Position.Add(component.Velocity.Scale(float64(delta.Seconds())))
}
