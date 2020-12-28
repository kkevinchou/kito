package steering

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/types"
)

type Seeker interface {
	types.Positionable
	Velocity() mgl64.Vec3
	Mass() float64
	MaxSpeed() float64
}

type SeekComponent struct {
	Entity Seeker
	target mgl64.Vec3
	active bool
}

func (s *SeekComponent) CalculateSteeringVelocity() mgl64.Vec3 {
	if !s.active {
		return mgl64.Vec3{}
	}

	desiredVelocity := s.target.Sub(s.Entity.Position()).Normalize().Mul(s.Entity.MaxSpeed())
	return desiredVelocity.Sub(s.Entity.Velocity()).Mul(1.0 / s.Entity.Mass())
}

func (s *SeekComponent) SetTarget(v mgl64.Vec3) {
	s.active = true
	s.target = v
}

func (s *SeekComponent) SetSeekActive(active bool) {
	s.active = active
}
