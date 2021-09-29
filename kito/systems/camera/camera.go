package camera

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/singleton"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	mouseWheelSensitivity float64 = 0.5
	zoomDecay             float64 = 0.1
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
}

type CameraSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
	return &s
}

func (s *CameraSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.CameraComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CameraSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()

	for _, camera := range s.entities {
		playerID := camera.GetComponentContainer().FollowComponent.FollowTargetEntityID
		handleCameraControls(delta, camera, s.world, singleton.PlayerInput[playerID])
	}
}

func handleCameraControls(delta time.Duration, entity entities.Entity, world World, frameInput input.Input) {
	cc := entity.GetComponentContainer()
	followComponent := cc.FollowComponent
	transformComponent := cc.TransformComponent

	if followComponent == nil {
		return
	}

	var xRel, yRel float64

	keyboardInput := frameInput.KeyboardInput
	mouseInput := frameInput.MouseInput

	var mouseSensitivity float64 = 0.005
	if mouseInput.LeftButtonDown && !mouseInput.MouseMotionEvent.IsZero() {
		xRel += -mouseInput.MouseMotionEvent.XRel * mouseSensitivity
		yRel += -mouseInput.MouseMotionEvent.YRel * mouseSensitivity
	}

	// handle camera controls with arrow keys
	var keyboardSensitivity float64 = 0.01
	if key, ok := keyboardInput[input.KeyboardKeyRight]; ok && key.Event == input.KeyboardEventDown {
		xRel += keyboardSensitivity
	}
	if key, ok := keyboardInput[input.KeyboardKeyLeft]; ok && key.Event == input.KeyboardEventDown {
		xRel += -keyboardSensitivity
	}
	if key, ok := keyboardInput[input.KeyboardKeyUp]; ok && key.Event == input.KeyboardEventDown {
		yRel += -keyboardSensitivity
	}
	if key, ok := keyboardInput[input.KeyboardKeyDown]; ok && key.Event == input.KeyboardEventDown {
		yRel += keyboardSensitivity
	}

	forwardVector := transformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
	upVector := transformComponent.Orientation.Rotate(mgl64.Vec3{0, 1, 0})
	// there's probably away to get the right vector directly rather than going crossing the up vector :D
	rightVector := forwardVector.Cross(upVector)

	// calculate the quaternion for the delta in rotation
	deltaRotationX := mgl64.QuatRotate(yRel, rightVector)         // pitch
	deltaRotationY := mgl64.QuatRotate(xRel, mgl64.Vec3{0, 1, 0}) // yaw
	deltaRotation := deltaRotationY.Mul(deltaRotationX)

	newOrientation := deltaRotation.Mul(transformComponent.Orientation)

	// don't let the camera go upside down
	if newOrientation.Rotate(mgl64.Vec3{0, 1, 0})[1] < 0 {
		newOrientation = transformComponent.Orientation
	}

	zoomDirection := libutils.NormalizeF64(followComponent.ZoomSpeed)

	if mouseInput.MouseWheelDelta != 0 {
		// allow the buildup of zoom velocity if we have continuous mouse wheeling
		// allow instantaneous direction change
		if !libutils.SameSign(zoomDirection, float64(mouseInput.MouseWheelDelta)) {
			followComponent.ZoomSpeed = 0
		}

		zoomDirection = libutils.NormalizeF64(float64(mouseInput.MouseWheelDelta))
		followComponent.ZoomSpeed += float64(mouseInput.MouseWheelDelta)
	}

	// decay zoom velocity
	followComponent.ZoomSpeed -= zoomDecay * zoomDirection
	if !libutils.SameSign(zoomDirection, followComponent.ZoomSpeed) {
		followComponent.ZoomSpeed = 0
	}

	followComponent.Zoom += followComponent.ZoomSpeed * mouseWheelSensitivity

	if followComponent.FollowDistance-followComponent.Zoom >= followComponent.MaxFollowDistance {
		followComponent.Zoom = -(followComponent.MaxFollowDistance - followComponent.FollowDistance)
	} else if followComponent.FollowDistance-followComponent.Zoom < 5 {
		followComponent.Zoom = followComponent.FollowDistance - 5
	}

	target, err := world.GetEntityByID(followComponent.FollowTargetEntityID)
	if err != nil {
		fmt.Println("failed to find target entity with ID", followComponent.FollowTargetEntityID)
		return
	}
	targetComponentContainer := target.GetComponentContainer()
	targetPosition := targetComponentContainer.TransformComponent.Position.Add(mgl64.Vec3{0, followComponent.YOffset, 0})
	transformComponent.Position = newOrientation.Rotate(mgl64.Vec3{0, 0, followComponent.FollowDistance - followComponent.Zoom}).Add(targetPosition)
	transformComponent.Orientation = newOrientation
}