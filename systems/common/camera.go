package common

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/input"
)

func HandleFollowCameraControls(entity entities.Entity, world World, frameInput input.Input) {
	followComponent := entity.GetComponentContainer().FollowComponent
	transformComponent := entity.GetComponentContainer().TransformComponent

	if followComponent == nil {
		return
	}

	entity, err := world.GetEntityByID(followComponent.FollowTargetEntityID)
	if err != nil {
		fmt.Println("failed to find target entity with ID", followComponent.FollowTargetEntityID)
		return
	}

	targetComponentContainer := entity.GetComponentContainer()
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

	targetToCamera := transformComponent.Position.Sub(targetPosition)
	transformComponent.Position = targetPosition.Add(deltaRotation.Rotate(targetToCamera).Normalize().Mul(followComponent.FollowDistance))
	transformComponent.Orientation = nextOrientation
	transformComponent.UpVector = deltaRotation.Rotate(transformComponent.UpVector)
}
