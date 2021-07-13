package common

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/types"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

func UpdateCharacterController(entity entities.Entity, world World, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	physicsComponent := componentContainer.PhysicsComponent

	keyboardInput := frameInput.KeyboardInput

	controlVector := GetControlVector(keyboardInput)

	forwardVector := mgl64.Vec3{0, 0, -1}
	rightVector := mgl64.Vec3{1, 0, 0}

	if tpcComponent := componentContainer.ThirdPersonControllerComponent; tpcComponent != nil {
		camera, err := world.GetEntityByID(tpcComponent.CameraID)
		if err != nil {
			panic(err)
		}
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
	var moveSpeed float64 = 20

	if !utils.Vec3IsZero(movementVector) {
		normalizedMovementVector := movementVector.Normalize()
		impulse := types.Impulse{
			Vector:    normalizedMovementVector.Mul(moveSpeed),
			DecayRate: 5,
		}
		physicsComponent.ApplyImpulse("controllerMove", impulse)
	}

	if controlVector.Y() > 0 {
		impulse := types.Impulse{
			Vector:    mgl64.Vec3{0, 40, 0},
			DecayRate: 1,
		}
		physicsComponent.ApplyImpulse("jumper", impulse)
	}
}
