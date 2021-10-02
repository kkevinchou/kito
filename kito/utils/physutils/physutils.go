package physutils

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	fullDecayThreshold = float64(0.05)
	gravity            = float64(60)
)

func PhysicsStep(delta time.Duration, entities []entities.Entity, playerID int) {
	accelerationDueToGravity := mgl64.Vec3{0, -gravity, 0}

	for _, entity := range entities {
		// TODO: I don't like this here, probably a more elegant solution
		if utils.IsClient() && entity.GetID() != playerID {
			continue
		}

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

		// updating orientation along velocity
		velocityWithoutY := mgl64.Vec3{velocity[0], 0, velocity[2]}
		if !libutils.Vec3IsZero(velocityWithoutY) {
			// Note, this will bug out if we look directly up or directly down. This
			// is due to issues looking at objects that are along our "up" vector.
			// I believe this is due to us losing sense of what a "right" vector is.
			// This code will likely change when we do animation blending in the animator
			transformComponent.Orientation = libutils.QuatLookAt(mgl64.Vec3{0, 0, 0}, velocityWithoutY.Normalize(), mgl64.Vec3{0, 1, 0})
		}

		// if entity.GetID() == 70000 {
		// 	fmt.Printf("[CF:%d] POST PHYSICS %v\n", world.GetSingleton().CommandFrame, transformComponent.Position)
		// }
	}
}
