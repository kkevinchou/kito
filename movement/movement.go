package movement

import (
	"time"

	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/physics"
)

type Moveable interface {
	physics.PhysicsI
	CalculateSteeringVelocity() vector.Vector
}

type MovementSystem struct {
	moveables []Moveable
}

func (m *MovementSystem) Register(moveable Moveable) {
	m.moveables = append(m.moveables, moveable)
}

func NewMovementSystem() MovementSystem {
	m := MovementSystem{}
	m.moveables = make([]Moveable, 0)
	return m
}

func (m *MovementSystem) Update(delta time.Duration) {
	for _, moveable := range m.moveables {
		physComp := moveable.GetPhysicsComponent()
		steeringVelocity := moveable.CalculateSteeringVelocity()
		physComp.Velocity = physComp.Velocity.Add(steeringVelocity).Clamp(physComp.MaxSpeed)
		moveable.Update(delta)
	}
}
