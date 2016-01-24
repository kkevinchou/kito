package movement

import (
	"time"

	"github.com/kkevinchou/ant/lib/math/vector"
)

type Moveable interface {
	Update(time.Duration)
	Velocity() vector.Vector
	SetVelocity(vector.Vector)
	MaxSpeed() float64
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
		steeringVelocity := moveable.CalculateSteeringVelocity()
		moveable.SetVelocity(moveable.Velocity().Add(steeringVelocity).Clamp(moveable.MaxSpeed()))
		moveable.Update(delta)
	}
}
