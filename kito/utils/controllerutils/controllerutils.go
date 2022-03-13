package controllerutils

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/input"
)

const (
	gravity   float64 = 2
	jumpSpeed float64 = 3
)

var (
	accelerationDueToGravity = mgl64.Vec3{0, -gravity, 0}
)

func UpdateCharacterController(delta time.Duration, entity entities.Entity, camera entities.Entity, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	cameraComponentContainer := camera.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent
	tpcComponent := componentContainer.ThirdPersonControllerComponent

	keyboardInput := frameInput.KeyboardInput
	controlVector := common.GetControlVector(keyboardInput)

	forwardVector := mgl64.Vec3{0, 0, -1}
	rightVector := mgl64.Vec3{1, 0, 0}

	forwardVector = cameraComponentContainer.TransformComponent.Orientation.Rotate(forwardVector)
	forwardVector[1] = 0
	forwardVector = forwardVector.Normalize()

	rightVector = cameraComponentContainer.TransformComponent.Orientation.Rotate(rightVector)
	rightVector[1] = 0
	rightVector.Normalize()

	forwardVector = forwardVector.Mul(controlVector.Z())
	rightVector = rightVector.Mul(controlVector.X())

	if tpcComponent.Grounded {
		tpcComponent.Grounded = false
		jumpVector := mgl64.Vec3{0, 1, 0}.Mul(controlVector.Y() * jumpSpeed)
		tpcComponent.Velocity = tpcComponent.Velocity.Add(jumpVector)
	}
	tpcComponent.Velocity = tpcComponent.Velocity.Add(accelerationDueToGravity.Mul(delta.Seconds()))

	movementVector := forwardVector.Add(rightVector)
	componentContainer.ThirdPersonControllerComponent.MovementVector = movementVector
	transformComponent.Position = transformComponent.Position.Add(movementVector).Add(tpcComponent.Velocity)
}

func ResolveControllerCollision(entity entities.Entity) {
	cc := entity.GetComponentContainer()
	colliderComponent := cc.ColliderComponent
	transformComponent := cc.TransformComponent
	tpcComponent := cc.ThirdPersonControllerComponent
	contactManifolds := colliderComponent.ContactManifolds
	if contactManifolds != nil {
		separatingVector := combineSeparatingVectors(contactManifolds)
		transformComponent.Position = transformComponent.Position.Add(separatingVector)
		tpcComponent.Grounded = true
		tpcComponent.Velocity = mgl64.Vec3{}
	} else {
		// no collisions were detected (i.e. the ground)
		// physicsComponent.Grounded = false
	}
}

func combineSeparatingVectors(contactManifolds []*collision.ContactManifold) mgl64.Vec3 {
	// only add separating vectors which we haven't seen before. ideally
	// this should handle cases where separating vectors are a basis of another
	// and avoid "overcounting" separating vectors
	seenSeparatingVectors := []mgl64.Vec3{}
	var separatingVector mgl64.Vec3
	for _, contactManifold := range contactManifolds {
		for _, contact := range contactManifold.Contacts {
			seen := false
			for _, v := range seenSeparatingVectors {
				if contact.SeparatingVector.ApproxEqual(v) {
					seen = true
					break
				}
			}
			if !seen {
				separatingVector = separatingVector.Add(contact.SeparatingVector)
				seenSeparatingVectors = append(seenSeparatingVectors, contact.SeparatingVector)
			}
		}
	}
	return separatingVector
}
