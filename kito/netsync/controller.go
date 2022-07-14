package netsync

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	jumpSpeed      float64 = 150
	zipSpeed       float64 = 400
	equalThreshold float64 = 1e-5

	// a value of 1 means the normal vector of what you're on must be exactly Vec3{0, 1, 0}
	groundedStrictness = 0.85

	// the maximum number of times a distinct entity can have their collision resolved
	// this presents the collision resolution phase to go on forever
	resolveCountMax = 10
)

// BaseVelocity - does not involve controller velocities (e.g. WASD)
// Velocity - actual observable velocity by external systems that includes movement velocities (e.g. WASD)
//          - computed each frame
func UpdateCharacterController(delta time.Duration, entity entities.Entity, camera entities.Entity, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent
	tpcComponent := componentContainer.ThirdPersonControllerComponent

	keyboardInput := frameInput.KeyboardInput
	controlVector := getControlVector(keyboardInput)

	// handle casting
	if _, ok := keyboardInput[input.KeyboardKeyQ]; ok {
		notepad := componentContainer.NotepadComponent
		notepad.LastAction = components.ActionCast
	}

	// handle jumping
	if controlVector.Y() > 0 && tpcComponent.Grounded {
		jumpVelocity := mgl64.Vec3{0, 1, 0}.Mul(controlVector.Y() * jumpSpeed)
		tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(jumpVelocity)
		tpcComponent.Grounded = false
	}

	// handle zip movement
	if _, ok := keyboardInput[input.KeyboardKeyE]; ok {
		if !libutils.Vec3ApproxEqualZero(tpcComponent.ZipVelocity) {
			tpcComponent.ZipVelocity = tpcComponent.ZipVelocity.Normalize().Mul(zipSpeed)
		} else {
			cameraView := frameInput.CameraOrientation.Rotate(mgl64.Vec3{0, 1, -5})
			tpcComponent.ZipVelocity = cameraView.Normalize().Mul(zipSpeed)
		}
	} else {
		tpcComponent.ZipVelocity = tpcComponent.ZipVelocity.Mul(.99)
		if libutils.Vec3ApproxEqualZero(tpcComponent.ZipVelocity) {
			tpcComponent.ZipVelocity = mgl64.Vec3{}
		}
	}

	// handle controller movement
	movementDir := calculateMovementDir(frameInput.CameraOrientation, controlVector)
	tpcComponent.MovementSpeed = computeMoveSpeed(tpcComponent.MovementSpeed)
	tpcComponent.ControllerVelocity = movementDir.Mul(tpcComponent.MovementSpeed)

	// apply all the various velocity adjustments
	tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(settings.AccelerationDueToGravity.Mul(delta.Seconds()))
	tpcComponent.Velocity = tpcComponent.BaseVelocity
	tpcComponent.Velocity = tpcComponent.Velocity.Add(tpcComponent.ControllerVelocity)
	tpcComponent.Velocity = tpcComponent.Velocity.Add(tpcComponent.ZipVelocity)

	transformComponent.Position = transformComponent.Position.Add(tpcComponent.Velocity.Mul(delta.Seconds()))

	// safeguard falling off the map
	if transformComponent.Position[1] < -1000 {
		transformComponent.Position[1] = 25
	}

	if !libutils.Vec3ApproxEqualZero(tpcComponent.ControllerVelocity) {
		transformComponent.Orientation = libutils.QuatLookAt(mgl64.Vec3{0, 0, 0}, tpcComponent.ControllerVelocity.Normalize(), mgl64.Vec3{0, 1, 0})
	} else {
		tpcComponent.MovementSpeed = 0
	}
}

func minSeparatingVector(contacts []*collision.Contact) mgl64.Vec3 {
	minVector := contacts[0].SeparatingVector
	minDistance := contacts[0].SeparatingDistance

	// one manifold for each object that's being collided with
	for _, contact := range contacts {
		if contact.SeparatingDistance < minDistance {
			minVector = contact.SeparatingVector
			minDistance = contact.SeparatingDistance
		}
	}

	return minVector
}

func hashVec(v mgl64.Vec3) string {
	return fmt.Sprintf("%.4f %.4f %.4f", v[0], v[1], v[2])
}

// combineSeparatingVectors: for triangles that have the same normal, only use the largest separating vector
// TODO: there might be some edge cases I'm not considering but visually it looks good enough
// there is still some jittering that happens when we have multiple triangles of varying
// normals. We probably need to be able to "merge" those somehow rather than naively summing
// them together like we are doing right now
func combineSeparatingVectors(contacts []*collision.Contact) mgl64.Vec3 {
	seenNormals := map[string]mgl64.Vec3{}
	for _, contact := range contacts {
		normalHash := hashVec(contact.Normal)
		if v, ok := seenNormals[normalHash]; ok {
			if contact.SeparatingVector.LenSqr() > v.LenSqr() {
				seenNormals[normalHash] = contact.SeparatingVector
			}
			continue
		}

		seenNormals[normalHash] = contact.SeparatingVector
	}

	var finalSeparatingVector mgl64.Vec3
	for _, v := range seenNormals {
		finalSeparatingVector = finalSeparatingVector.Add(v)
	}
	return finalSeparatingVector
}

func computeMoveSpeed(movementSpeed float64) float64 {
	if movementSpeed < 60 {
		return movementSpeed + 15
	} else if movementSpeed < 100 {
		return movementSpeed + 2
	}
	return movementSpeed
}

// movementDir does not include Y values
func calculateMovementDir(cameraOrientation mgl64.Quat, controlVector mgl64.Vec3) mgl64.Vec3 {
	forwardVector := cameraOrientation.Rotate(mgl64.Vec3{0, 0, -1})
	forwardVector = forwardVector.Normalize().Mul(controlVector.Z())
	forwardVector[1] = 0

	rightVector := cameraOrientation.Rotate(mgl64.Vec3{1, 0, 0})
	rightVector = rightVector.Normalize().Mul(controlVector.X())
	rightVector[1] = 0

	movementDir := forwardVector.Add(rightVector)
	if movementDir.LenSqr() > 0 {
		return movementDir.Normalize()
	}

	return movementDir
}

func calculateCameraForwardRightVec(camera entities.Entity) (mgl64.Vec3, mgl64.Vec3) {
	cc := camera.GetComponentContainer()
	forwardVector := cc.TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1}).Normalize()
	rightVector := cc.TransformComponent.Orientation.Rotate(mgl64.Vec3{1, 0, 0}).Normalize()

	return forwardVector, rightVector
}

func getControlVector(keyboardInput input.KeyboardInput) mgl64.Vec3 {
	var controlVector mgl64.Vec3
	if key, ok := keyboardInput[input.KeyboardKeyW]; ok && key.Event == input.KeyboardEventDown {
		controlVector[2]++
	}
	if key, ok := keyboardInput[input.KeyboardKeyS]; ok && key.Event == input.KeyboardEventDown {
		controlVector[2]--
	}
	if key, ok := keyboardInput[input.KeyboardKeyA]; ok && key.Event == input.KeyboardEventDown {
		controlVector[0]--
	}
	if key, ok := keyboardInput[input.KeyboardKeyD]; ok && key.Event == input.KeyboardEventDown {
		controlVector[0]++
	}
	if key, ok := keyboardInput[input.KeyboardKeyLShift]; ok && key.Event == input.KeyboardEventDown {
		controlVector[1]--
	}
	if key, ok := keyboardInput[input.KeyboardKeySpace]; ok && key.Event == input.KeyboardEventDown {
		controlVector[1]++
	}
	return controlVector
}