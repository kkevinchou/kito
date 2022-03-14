package controllerutils

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	gravity   float64 = 250
	jumpSpeed float64 = 150
	moveSpeed float64 = 100
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

	forwardVector := cameraComponentContainer.TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
	forwardVector[1] = 0
	forwardVector = forwardVector.Normalize()

	rightVector := cameraComponentContainer.TransformComponent.Orientation.Rotate(mgl64.Vec3{1, 0, 0})
	rightVector[1] = 0
	rightVector.Normalize()

	forwardVector = forwardVector.Mul(controlVector.Z())
	rightVector = rightVector.Mul(controlVector.X())

	if tpcComponent.Grounded {
		jumpVelocity := mgl64.Vec3{0, 1, 0}.Mul(controlVector.Y() * jumpSpeed)
		tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(jumpVelocity)
		tpcComponent.Grounded = false
	}
	tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(accelerationDueToGravity.Mul(delta.Seconds()))
	movementVelocity := forwardVector.Add(rightVector).Mul(moveSpeed)
	tpcComponent.Velocity = tpcComponent.BaseVelocity.Add(movementVelocity)
	transformComponent.Position = transformComponent.Position.Add(tpcComponent.Velocity.Mul(delta.Seconds()))
	if transformComponent.Position[1] < -1000 {
		transformComponent.Position[1] = 25
	}

	if !libutils.Vec3IsZero(movementVelocity) {
		transformComponent.Orientation = libutils.QuatLookAt(mgl64.Vec3{0, 0, 0}, movementVelocity.Normalize(), mgl64.Vec3{0, 1, 0})
	}
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
		tpcComponent.BaseVelocity = mgl64.Vec3{}
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
