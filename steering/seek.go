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
}

func (s *SeekComponent) CalculateSteeringVelocity() vector.Vector {
	desiredVelocity := s.target.Sub(s.Entity.Position()).Normalize().Scale(s.Entity.MaxSpeed())
	return desiredVelocity.Sub(s.Entity.Velocity()).Scale(1.0 / s.Entity.Mass())
}

func (s *SeekComponent) SetTarget(v vector.Vector) {
	s.target = v
}
