package physics

import (
	"time"

	"github.com/kkevinchou/kito/lib/utils"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
)

const (
	fullDecayThreshold = float64(0.05)
)

type World interface{}

type PhysicsSystem struct {
	world    World
	entities []entities.Entity
}

func NewPhysicsSystem(world World) *PhysicsSystem {
	return &PhysicsSystem{
		world:    world,
		entities: []entities.Entity{},
	}
}

func (s *PhysicsSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *PhysicsSystem) Update(delta time.Duration) {
	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		physicsComponent := componentContainer.PhysicsComponent
		transformComponent := componentContainer.TransformComponent

		var totalImpulse mgl64.Vec3
		for name := range physicsComponent.Impulses {
			impulse := physicsComponent.Impulses[name]
			impulse.ElapsedTime = impulse.ElapsedTime + delta
			decayRatio := 1.0 - (impulse.ElapsedTime.Seconds() * impulse.DecayRate)
			if decayRatio < 0 {
				decayRatio = 0
			}

			if decayRatio < fullDecayThreshold {
				delete(physicsComponent.Impulses, name)
			} else {
				realImpulse := impulse.Vector.Mul(decayRatio)
				totalImpulse = totalImpulse.Add(realImpulse)
			}
		}

		velocity := physicsComponent.Velocity.Add(totalImpulse)
		newPos := transformComponent.Position.Add(velocity.Mul(delta.Seconds()))
		transformComponent.Position = newPos

		if !utils.Vec3IsZero(velocity) {
			transformComponent.ForwardVector = velocity.Normalize()
			// Note, this will bug out if we look directly up or directly down. This
			// is due to issues looking at objects that are along our "up" vector.
			// I believe this is due to us losing sense of what a "right" vector is.
			// This code will likely change when we do animation blending in the animator
			transformComponent.ViewQuaternion = utils.QuatLookAt(mgl64.Vec3{0, 0, 0}, transformComponent.ForwardVector, mgl64.Vec3{0, 1, 0})
		}
	}
}
