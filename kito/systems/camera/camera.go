package camera

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/singleton"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	farMouseWheelSensitivity  float64 = 2.5
	nearMouseWheelSensitivity float64 = 1.5
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) entities.Entity
	QueryEntity(componentFlags int) []entities.Entity
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

func (s *CameraSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()

	for _, camera := range s.world.QueryEntity(components.ComponentFlagCamera | components.ComponentFlagControl) {
		playerID := camera.GetComponentContainer().ControlComponent.PlayerID
		newOrientation := handleCameraControls(delta, camera, s.world, singleton.PlayerInput[playerID])
		currentInput := singleton.PlayerInput[playerID]
		currentInput.CameraOrientation = newOrientation
		singleton.PlayerInput[playerID] = currentInput
	}
}

func handleCameraControls(delta time.Duration, entity entities.Entity, world World, frameInput input.Input) mgl64.Quat {
	cc := entity.GetComponentContainer()
	cameraComponent := cc.CameraComponent
	transformComponent := cc.TransformComponent

	var xRel, yRel float64

	keyboardInput := frameInput.KeyboardInput
	mouseInput := frameInput.MouseInput

	var mouseSensitivity float64 = 0.005
	if mouseInput.Buttons[1] && !mouseInput.MouseMotionEvent.IsZero() {
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

	mouseWheelSensitivity := farMouseWheelSensitivity
	if cameraComponent.FollowDistance < 50 {
		mouseWheelSensitivity = nearMouseWheelSensitivity
	}

	if mouseInput.MouseWheelDelta != 0 {
		currentMouseZoomDirection := libutils.NormalizeF64(float64(mouseInput.MouseWheelDelta))
		cameraComponent.ZoomSpeed = currentMouseZoomDirection * -mouseWheelSensitivity
	}

	// decay zoom velocity
	cameraComponent.ZoomSpeed *= 0.90
	if math.Abs(cameraComponent.ZoomSpeed) < 0.01 {
		cameraComponent.ZoomSpeed = 0
	}

	cameraComponent.FollowDistance += cameraComponent.ZoomSpeed

	if cameraComponent.FollowDistance >= cameraComponent.MaxFollowDistance {
		cameraComponent.FollowDistance = cameraComponent.MaxFollowDistance
	} else if cameraComponent.FollowDistance < 5 {
		cameraComponent.FollowDistance = 5
	}

	target := world.GetEntityByID(cameraComponent.FollowTargetEntityID)
	if target == nil {
		fmt.Println("failed to find target entity with ID", cameraComponent.FollowTargetEntityID)
		return mgl64.QuatIdent()
	}
	targetComponentContainer := target.GetComponentContainer()
	targetPosition := targetComponentContainer.TransformComponent.Position.Add(mgl64.Vec3{0, cameraComponent.YOffset, 0})
	transformComponent.Position = newOrientation.Rotate(mgl64.Vec3{0, 0, cameraComponent.FollowDistance}).Add(targetPosition)
	transformComponent.Orientation = newOrientation
	return newOrientation
}
