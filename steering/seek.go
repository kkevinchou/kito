package steering

import "github.com/kkevinchou/ant/math/vector"

type Seekable interface {
	Position() vector.Vector
	Velocity() vector.Vector
	Mass() float64
	MaxSpeed() float64
}

type SeekComponent struct {
	Entity Seekable
	target vector.Vector
	active bool
}

func (s *SeekComponent) CalculateSteeringVelocity() vector.Vector {
	if !s.active {
		return vector.Vector{}
	}

	desiredVelocity := s.target.Sub(s.Entity.Position()).Normalize().Scale(s.Entity.MaxSpeed())
	return desiredVelocity.Sub(s.Entity.Velocity()).Scale(1.0 / s.Entity.Mass())
}

func (s *SeekComponent) SetTarget(v vector.Vector) {
	s.active = true
	s.target = v
}

func (s *SeekComponent) SetSeekActive(active bool) {
	s.active = active
}
