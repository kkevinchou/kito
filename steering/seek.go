package steering

import (
	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/physics"
)

type SeekComponent struct {
	physicsComponent *physics.PhysicsComponent
	target           vector.Vector
}

func (s *SeekComponent) Initialize(p *physics.PhysicsComponent) {
	s.physicsComponent = p
}

func (s *SeekComponent) CalculateSteeringVelocity() vector.Vector {
	physComp := s.physicsComponent
	desiredVelocity := s.target.Sub(physComp.Position).Normalize().Scale(physComp.MaxSpeed)
	return desiredVelocity.Sub(physComp.Velocity).Scale(1.0 / physComp.Mass)
}

func (s *SeekComponent) SetTarget(v vector.Vector) {
	s.target = v
}
