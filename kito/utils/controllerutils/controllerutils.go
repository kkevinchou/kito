package controllerutils

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/lib/input"
)

func UpdateCharacterController(entity entities.Entity, camera entities.Entity, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent

	keyboardInput := frameInput.KeyboardInput
	controlVector := common.GetControlVector(keyboardInput)

	forwardVector := mgl64.Vec3{0, 0, -1}
	rightVector := mgl64.Vec3{1, 0, 0}
	upVector := mgl64.Vec3{0, 1, 0}

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
	upVector = upVector.Mul(controlVector.Y())

	movementVector := forwardVector.Add(rightVector).Add(upVector.Mul(5))

	componentContainer.ThirdPersonControllerComponent.MovementVector = movementVector
	transformComponent.Position = transformComponent.Position.Add(movementVector).Add(mgl64.Vec3{0, -1, 0})
}
