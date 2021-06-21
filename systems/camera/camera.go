package camera

import (
	"fmt"
	"math"
	"time"

	"github.com/kkevinchou/kito/components"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/systems/base"
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
	*base.BaseSystem
	world World
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
	return &s
}

func (s *CameraSystem) RegisterEntity(entity entities.Entity) {
}

func (s *CameraSystem) Update(delta time.Duration) {
	camera := s.world.GetCamera()
	componentContainer := camera.GetComponentContainer()
	s.handleFollowCameraControls(componentContainer)
}

// this might belong in some kind of movement or pathfinding system that handles "following" logic.
// putting this here for now until more than just cameras need to follow a target
func (s *CameraSystem) handleFollowCameraControls(componentContainer *components.ComponentContainer) {
	followComponent := componentContainer.FollowComponent
	transformComponent := componentContainer.TransformComponent

	if followComponent == nil {
		return
	}

	entity, err := s.world.GetEntityByID(followComponent.FollowTargetEntityID)
	if err != nil {
		fmt.Println("failed to find target entity with ID", followComponent.FollowTargetEntityID)
		return
	}

	targetComponentContainer := entity.GetComponentContainer()
	targetPosition := targetComponentContainer.TransformComponent.Position

	var xRel, yRel float64

	singleton := s.world.GetSingleton()
	mouseInput := singleton.GetMouseInput()

	if mouseInput != nil {
		var mouseSensitivity float64 = 0.005
		if mouseInput.LeftButtonDown && mouseInput.MouseMotionEvent != nil {
			xRel += -mouseInput.MouseMotionEvent.XRel * mouseSensitivity
			yRel += -mouseInput.MouseMotionEvent.YRel * mouseSensitivity
		}
	}

	// handle camera controls with arrow keys
	if singleton.GetKeyboardInputSet() != nil {
		keyboardInput := *singleton.GetKeyboardInputSet()
		var keyboardSensitivity float64 = 0.01
		if key, ok := keyboardInput[types.KeyboardKeyRight]; ok && key.Event == types.KeyboardEventDown {
			xRel += keyboardSensitivity
		}
		if key, ok := keyboardInput[types.KeyboardKeyLeft]; ok && key.Event == types.KeyboardEventDown {
			xRel += -keyboardSensitivity
		}
		if key, ok := keyboardInput[types.KeyboardKeyUp]; ok && key.Event == types.KeyboardEventDown {
			yRel += -keyboardSensitivity
		}
		if key, ok := keyboardInput[types.KeyboardKeyDown]; ok && key.Event == types.KeyboardEventDown {
			yRel += keyboardSensitivity
		}
	}

	forwardVector := transformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
	rightVector := forwardVector.Cross(transformComponent.UpVector)
	transformComponent.Position = targetPosition.Add(forwardVector.Mul(-1).Mul(followComponent.FollowDistance))

	// calculate the quaternion for the delta in rotation
	deltaRotationX := mgl64.QuatRotate(yRel, rightVector)         // pitch
	deltaRotationY := mgl64.QuatRotate(xRel, mgl64.Vec3{0, 1, 0}) // yaw
	deltaRotation := deltaRotationY.Mul(deltaRotationX)

	nextOrientation := deltaRotation.Mul(transformComponent.Orientation)
	nextForwardVector := nextOrientation.Rotate(mgl64.Vec3{0, 0, -1})

	// if we're nearly pointing directly downwards or upwards - stop camera movement
	// TODO: do this in a better way
	if mgl64.FloatEqualThreshold(math.Abs(nextForwardVector[1]), 1, 0.001) {
		return
	}

	targetToCamera := transformComponent.Position.Sub(targetPosition)
	transformComponent.Position = targetPosition.Add(deltaRotation.Rotate(targetToCamera).Normalize().Mul(followComponent.FollowDistance))
	transformComponent.Orientation = nextOrientation
	transformComponent.UpVector = deltaRotation.Rotate(transformComponent.UpVector)
}

// controlled cameras are cameras that can move independently
func (s *CameraSystem) handleFreeCamera(componentContainer *components.ComponentContainer) {
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

	impulse := types.Impulse{}
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
