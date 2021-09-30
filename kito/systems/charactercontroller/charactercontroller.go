package charactercontroller

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/input"
	libutils "github.com/kkevinchou/kito/lib/libutils"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type CharacterControllerSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCharacterControllerSystem(world World) *CharacterControllerSystem {
	return &CharacterControllerSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CharacterControllerSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil && componentContainer.ThirdPersonControllerComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CharacterControllerSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()
	singleton := s.world.GetSingleton()

	for _, player := range playerManager.GetPlayers() {
		entity, err := s.world.GetEntityByID(player.ID)
		if err != nil {
			continue
		}

		cameraID := entity.GetComponentContainer().ThirdPersonControllerComponent.CameraID
		camera, err := s.world.GetEntityByID(cameraID)
		if err != nil {
			continue
		}
		UpdateCharacterController(entity, camera, singleton.PlayerInput[player.ID])
	}
}

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
