package steering

import (
	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/physics"
)

type SeekComponent struct {
	Entity physics.PhysicsI
	target vector.Vector
}

func (s *SeekComponent) CalculateSteeringVelocity() vector.Vector {
	physComp := s.Entity.GetPhysicsComponent()
	desiredVelocity := s.target.Sub(physComp.Position).Normalize().Scale(physComp.MaxSpeed)
	return desiredVelocity.Sub(physComp.Velocity).Scale(1.0 / physComp.Mass)
}

func (s *SeekComponent) SetTarget(v vector.Vector) {
	s.target = v
}
