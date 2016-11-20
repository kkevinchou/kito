package movement

import (
	"time"

	"github.com/kkevinchou/ant/lib/math/vector"
)

type Mover interface {
	Update(time.Duration)
	Velocity() vector.Vector
	SetVelocity(vector.Vector)
	MaxSpeed() float64
	CalculateSteeringVelocity() vector.Vector
}

type MovementSystem struct {
	movers []Mover
}

func (m *MovementSystem) Register(mover Mover) {
	m.movers = append(m.movers, mover)
}

func NewMovementSystem() *MovementSystem {
	m := MovementSystem{}
	m.movers = make([]Mover, 0)
	return &m
}

func (m *MovementSystem) Update(delta time.Duration) {
	for _, mover := range m.movers {
		steeringVelocity := mover.CalculateSteeringVelocity()
		mover.SetVelocity(mover.Velocity().Add(steeringVelocity).Clamp(mover.MaxSpeed()))
		mover.Update(delta)
	}
}
