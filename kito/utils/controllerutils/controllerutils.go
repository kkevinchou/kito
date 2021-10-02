package controllerutils

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

func UpdateCharacterController(entity entities.Entity, camera entities.Entity, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	physicsComponent := componentContainer.PhysicsComponent

	keyboardInput := frameInput.KeyboardInput

	controlVector := common.GetControlVector(keyboardInput)

	forwardVector := mgl64.Vec3{0, 0, -1}
	rightVector := mgl64.Vec3{1, 0, 0}

	if tpcComponent := componentContainer.ThirdPersonControllerComponent; tpcComponent != nil {
		cameraComponentContainer := camera.GetComponentContainer()

		forwardVector = cameraComponentContainer.TransformComponent.Orientation.Rotate(forwardVector)
		forwardVector[1] = 0
		forwardVector = forwardVector.Normalize()

		rightVector = cameraComponentContainer.TransformComponent.Orientation.Rotate(rightVector)
		rightVector[1] = 0
		rightVector.Normalize()
	}

	forwardVector = forwardVector.Mul(controlVector.Z())
	rightVector = rightVector.Mul(controlVector.X())
	movementVector := forwardVector.Add(rightVector)
	var moveSpeed float64 = 100

	if !libutils.Vec3IsZero(movementVector) {
		normalizedMovementVector := movementVector.Normalize()
		impulse := types.Impulse{
			Vector:    normalizedMovementVector.Mul(moveSpeed),
			DecayRate: 5,
		}
		physicsComponent.ApplyImpulse("controllerMove", impulse)
	}

	if controlVector.Y() > 0 {
		impulse := types.Impulse{
			Vector:    mgl64.Vec3{0, 100, 0},
			DecayRate: 1,
		}
		physicsComponent.ApplyImpulse("jumper", impulse)
	}
}
