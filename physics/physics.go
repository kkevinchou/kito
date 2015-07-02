package physics

import (
	// "fmt"
	"github.com/kkevinchou/ant/math/vector"
	"time"
)

type PhysicsComposed interface {
	GetPhysicsComponent() *PhysicsComponent
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
