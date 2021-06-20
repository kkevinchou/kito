package physics

import (
	"time"

	"github.com/kkevinchou/kito/lib/utils"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
)

const (
	fullDecayThreshold = float64(0.05)
	gravity            = float64(60)
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
	accelerationDueToGravity := mgl64.Vec3{0, -gravity, 0}

	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		physicsComponent := componentContainer.PhysicsComponent
		transformComponent := componentContainer.TransformComponent

		// calculate impulses and their decay, this is meant for controller
		// actions that can "overwite" impulses
		var totalImpulse mgl64.Vec3
		for name, impulse := range physicsComponent.Impulses {
			decayRatio := 1.0 - (impulse.ElapsedTime.Seconds() * impulse.DecayRate)
			if decayRatio < 0 {
				decayRatio = 0
			}

			if decayRatio < fullDecayThreshold {
				delete(physicsComponent.Impulses, name)
				continue
			} else {
				realImpulse := impulse.Vector.Mul(decayRatio)
				totalImpulse = totalImpulse.Add(realImpulse)
			}

			// update the impulse
			impulse.ElapsedTime = impulse.ElapsedTime + delta
			physicsComponent.Impulses[name] = impulse
		}

		// calculate velocity adjusted by acceleration
		totalAcceleration := accelerationDueToGravity
		physicsComponent.Velocity = physicsComponent.Velocity.Add(totalAcceleration.Mul(delta.Seconds()))

		velocity := physicsComponent.Velocity.Add(totalImpulse)
		newPos := transformComponent.Position.Add(velocity.Mul(delta.Seconds()))

		// temporary hack to not fall through the ground
		if newPos[1] < 0 {
			newPos[1] = 0
			velocity[1] = 0
			physicsComponent.Velocity[1] = 0
		}

		transformComponent.Position = newPos

		// updating forward vector
		velocityWithoutY := mgl64.Vec3{velocity[0], 0, velocity[2]}
		if !utils.Vec3IsZero(velocityWithoutY) {
			transformComponent.ForwardVector = velocityWithoutY.Normalize()
			// Note, this will bug out if we look directly up or directly down. This
			// is due to issues looking at objects that are along our "up" vector.
			// I believe this is due to us losing sense of what a "right" vector is.
			// This code will likely change when we do animation blending in the animator
			transformComponent.Orientation = utils.QuatLookAt(mgl64.Vec3{0, 0, 0}, transformComponent.ForwardVector, mgl64.Vec3{0, 1, 0})
		}
	}
}
