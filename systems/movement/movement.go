package movement

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/types"
)

type Mover interface {
	Update(time.Duration)
	Velocity() mgl64.Vec3
	SetVelocity(mgl64.Vec3)
	MaxSpeed() float64
	CalculateSteeringVelocity() mgl64.Vec3
	MovementType() types.MovementType
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
		if mover.MovementType() == types.MovementTypeSteering {
			steeringVelocity := mover.CalculateSteeringVelocity()
			mover.SetVelocity(mover.Velocity().Add(steeringVelocity))
			if mover.Velocity().Len() > mover.MaxSpeed() {
				mover.SetVelocity(mover.Velocity().Normalize().Mul(mover.MaxSpeed()))
			}
			mover.Update(delta)
		}
	}
}
