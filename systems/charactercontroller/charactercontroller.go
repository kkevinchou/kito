package charactercontroller

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/systems/sysutils"
	"github.com/kkevinchou/kito/types"
)

type World interface {
	GetSingleton() types.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type CharacterControllerSystem struct {
	world    World
	entities []entities.Entity
}

func NewCharacterControllerSystem(world World) *CharacterControllerSystem {
	return &CharacterControllerSystem{world: world}
}

func (s *CharacterControllerSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil && componentContainer.ThirdPersonControllerComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CharacterControllerSystem) Update(delta time.Duration) {
	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		physicsComponent := componentContainer.PhysicsComponent

		singleton := s.world.GetSingleton()
		keyboardInput := *singleton.GetKeyboardInputSet()

		controlVector := sysutils.GetControlVector(keyboardInput)

		forwardVector := mgl64.Vec3{0, 0, -1}
		rightVector := mgl64.Vec3{1, 0, 0}

		if tpcComponent := componentContainer.ThirdPersonControllerComponent; tpcComponent != nil {
			camera, err := s.world.GetEntityByID(tpcComponent.CameraID)
			if err != nil {
				panic(err)
			}
			cameraComponentContainer := camera.GetComponentContainer()

			forwardVector = cameraComponentContainer.TransformComponent.ViewQuaternion.Rotate(forwardVector)
			// forwardVector = cameraComponentContainer.TransformComponent.ForwardVector
			forwardVector[1] = 0
			forwardVector = forwardVector.Normalize()

			rightVector = cameraComponentContainer.TransformComponent.ViewQuaternion.Rotate(rightVector)
			// rightVector = forwardVector.Cross(cameraComponentContainer.TransformComponent.UpVector)
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
				Vector:    mgl64.Vec3{0, 15, 0},
				DecayRate: 1,
			}
			physicsComponent.ApplyImpulse("jumper", impulse)
		}
	}
}
