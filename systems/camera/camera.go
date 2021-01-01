package camera

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/components"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/systems/sysutils"
	"github.com/kkevinchou/kito/types"
)

const (
	zoomSpeed float64 = 100
	moveSpeed float64 = 25
)

type Singleton interface {
	GetKeyboardInputSet() *types.KeyboardInput
}

type World interface {
	GetSingleton() types.Singleton
	GetCamera() entities.Entity
	GetEntityByID(id int) (entities.Entity, error)
}

type CameraSystem struct {
	world World
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		world: world,
	}
	return &s
}

func (s *CameraSystem) Update(delta time.Duration) {
	camera := s.world.GetCamera()

	componentContainer := camera.GetComponentContainer()
	controllerComponent := componentContainer.ControllerComponent

	if controllerComponent.Controlled {
		s.handleControlledCamera(componentContainer)
		return
	}

	if controllerComponent == nil || !controllerComponent.Controlled {
		s.handleUncontrolledCamera(componentContainer)
	}

}

// this might belong in some kind of movement or pathfinding system that handles "following" logic.
// putting this here for now until more than just cameras need to follow a target
func (s *CameraSystem) handleUncontrolledCamera(componentContainer *components.ComponentContainer) {
	followComponent := componentContainer.FollowComponent
	positionComponent := componentContainer.PositionComponent

	if followComponent == nil || followComponent.FollowTargetEntityID == nil {
		return
	}

	entity, err := s.world.GetEntityByID(*followComponent.FollowTargetEntityID)
	if err != nil {
		panic(err)
	}

	targetComponentContainer := entity.GetComponentContainer()

	targetPosition := targetComponentContainer.PositionComponent.Position
	positionComponent.Position = targetPosition
	positionComponent.Position[1] += 10
	positionComponent.Position[2] += 20
}

func (s *CameraSystem) handleControlledCamera(componentContainer *components.ComponentContainer) {
	physicsComponent := componentContainer.PhysicsComponent
	topDownViewComponent := componentContainer.TopDownViewComponent

	singleton := s.world.GetSingleton()
	keyboardInput := *singleton.GetKeyboardInputSet()
	controlVector := sysutils.GetControlVector(keyboardInput)

	zoomValue := 0
	if singleton.GetMouseInput() != nil {
		mouseInput := *singleton.GetMouseInput()
		if mouseInput.LeftButtonDown && mouseInput.MouseMotionEvent != nil {
			topDownViewComponent.UpdateView(mgl64.Vec2{mouseInput.MouseMotionEvent.XRel, mouseInput.MouseMotionEvent.YRel})
		}

		if mouseInput.MouseWheel == types.MouseWheelDirectionNeutral {
			zoomValue = 0
		} else if mouseInput.MouseWheel == types.MouseWheelDirectionUp {
			zoomValue = -1
		} else if mouseInput.MouseWheel == types.MouseWheelDirectionDown {
			zoomValue = 1
		} else {
			panic(fmt.Sprintf("unexpected mousewheel value %v", mouseInput.MouseWheel))
		}
	}

	if utils.Vec3IsZero(controlVector) && zoomValue == 0 {
		return
	}

	forwardVector := topDownViewComponent.Forward()
	zoomVector := forwardVector.Mul(float64(zoomValue))

	forwardVector = forwardVector.Mul(controlVector.Z())
	forwardVector[1] = 0

	rightVector := topDownViewComponent.Right()
	rightVector = rightVector.Mul(-controlVector.X())

	impulse := &types.Impulse{}
	if !utils.Vec3IsZero(controlVector) {
		impulse.Vector = forwardVector.Add(rightVector).Add(mgl64.Vec3{0, controlVector.Y(), 0}).Normalize().Mul(moveSpeed)
		impulse.DecayRate = 2.5
		physicsComponent.ApplyImpulse("cameraMove", impulse)
	}

	if zoomValue != 0 {
		impulse.Vector = zoomVector.Mul(zoomSpeed)
		impulse.DecayRate = 2.5
		physicsComponent.ApplyImpulse("cameraZoom", impulse)
	}
}
