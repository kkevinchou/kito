package camera

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities/singleton"

	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/utils"
	"github.com/kkevinchou/kito/systems/base"
)

const (
	moveSpeed             float64 = 25
	mouseWheelSensitivity float64 = 8
)

type Singleton interface {
	GetKeyboardInputSet() *input.KeyboardInput
}

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

	cc.EasingComponent.Update(delta)

	if followComponent == nil {
		return
	}

	target, err := world.GetEntityByID(followComponent.FollowTargetEntityID)
	if err != nil {
		fmt.Println("failed to find target entity with ID", followComponent.FollowTargetEntityID)
		return
	}

	targetComponentContainer := target.GetComponentContainer()
	targetPosition := targetComponentContainer.TransformComponent.Position

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

	// TBH this easing stuff feels awkward here, might belong in a separate system
	var easingValue float64
	if mouseInput.MouseWheelDelta != 0 {
		// allow the buildup of zoom velocity if we have continuous mouse wheeling
		if cc.EasingComponent.Active() && utils.SameSign(followComponent.ZoomDirection, mouseInput.MouseWheelDelta) {
			followComponent.ZoomVelocity += 1
		} else {
			followComponent.ZoomVelocity = 1
		}

		cc.EasingComponent.Start()
		if mouseInput.MouseWheelDelta > 0 {
			followComponent.ZoomDirection = 1
		} else {
			followComponent.ZoomDirection = -1
		}
	}

	if cc.EasingComponent.Active() {
		easingValue = cc.EasingComponent.GetValue()
		if easingValue == 1 {
			cc.EasingComponent.Stop()
			followComponent.ZoomDirection = 0
			followComponent.ZoomVelocity = 1
		}

	}

	followComponent.Zoom += math.Pow(1.8, followComponent.ZoomVelocity/2) * float64(followComponent.ZoomDirection) * (float64(1) / float64(delta.Milliseconds())) * mouseWheelSensitivity

	if followComponent.FollowDistance-followComponent.Zoom >= followComponent.MaxFollowDistance {
		followComponent.Zoom = -(followComponent.MaxFollowDistance - followComponent.FollowDistance)
	} else if followComponent.FollowDistance-followComponent.Zoom < 5 {
		followComponent.Zoom = followComponent.FollowDistance - 5
	}

	targetToCamera := transformComponent.Position.Sub(targetPosition).Normalize()
	transformComponent.Position = targetPosition.Add(deltaRotation.Rotate(targetToCamera).Mul(followComponent.FollowDistance - followComponent.Zoom))
	transformComponent.Orientation = nextOrientation
	transformComponent.UpVector = deltaRotation.Rotate(transformComponent.UpVector)
}
