package steering

import (
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/types"
)

type Seeker interface {
	types.Positionable
	Velocity() vector.Vector3
	Mass() float64
	MaxSpeed() float64
}

type SeekComponent struct {
	Entity Seeker
	target vector.Vector3
	active bool
}

func (s *SeekComponent) CalculateSteeringVelocity() vector.Vector3 {
	if !s.active {
		return vector.Vector3{}
	}

	desiredVelocity := s.target.Sub(s.Entity.Position()).Normalize().Scale(s.Entity.MaxSpeed())
	return desiredVelocity.Sub(s.Entity.Velocity()).Scale(1.0 / s.Entity.Mass())
}

func (s *SeekComponent) SetTarget(v vector.Vector3) {
	s.active = true
	s.target = v
}

func (s *SeekComponent) SetSeekActive(active bool) {
	s.active = active
}
