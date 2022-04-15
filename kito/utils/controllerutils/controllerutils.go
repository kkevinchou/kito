package controllerutils

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/common"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	gravity        float64 = 250
	jumpSpeed      float64 = 150
	zipSpeed       float64 = 400
	equalThreshold float64 = 1e-5
)

var (
	accelerationDueToGravity = mgl64.Vec3{0, -gravity, 0}
)

// BaseVelocity - does not involve controller velocities (e.g. WASD)
// Velocity - actual observable velocity by external systems that includes movement velocities (e.g. WASD)
//          - computed each frame
func UpdateCharacterController(delta time.Duration, entity entities.Entity, camera entities.Entity, frameInput input.Input) {
	componentContainer := entity.GetComponentContainer()
	transformComponent := componentContainer.TransformComponent
	tpcComponent := componentContainer.ThirdPersonControllerComponent

	keyboardInput := frameInput.KeyboardInput
	controlVector := common.GetControlVector(keyboardInput)

	// handle jumping
	if controlVector.Y() > 0 && tpcComponent.Grounded {
		jumpVelocity := mgl64.Vec3{0, 1, 0}.Mul(controlVector.Y() * jumpSpeed)
		tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(jumpVelocity)
		tpcComponent.Grounded = false
	}

	// handle zip movement
	if _, ok := keyboardInput[input.KeyboardKeyE]; ok {
		// forward, right := calculateCameraForwardRightVec(camera)
		if !libutils.Vec3ApproxEqualZero(tpcComponent.ZipVelocity) {
			tpcComponent.ZipVelocity = tpcComponent.ZipVelocity.Normalize().Mul(zipSpeed)
		} else {
			cameraView := camera.GetComponentContainer().TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
			tpcComponent.ZipVelocity = cameraView.Normalize().Mul(zipSpeed)
		}
	} else {
		tpcComponent.ZipVelocity = tpcComponent.ZipVelocity.Mul(.9)
		if libutils.Vec3ApproxEqualZero(tpcComponent.ZipVelocity) {
			tpcComponent.ZipVelocity = mgl64.Vec3{}
		}
	}

	// handle controller movement
	movementDir := calculateMovementDir(camera, controlVector)
	tpcComponent.MovementSpeed = computeMoveSpeed(tpcComponent.MovementSpeed)
	tpcComponent.ControllerVelocity = movementDir.Mul(tpcComponent.MovementSpeed)

	// apply all the various velocity adjustments
	tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(accelerationDueToGravity.Mul(delta.Seconds()))
	tpcComponent.Velocity = tpcComponent.BaseVelocity.Add(tpcComponent.ControllerVelocity)
	tpcComponent.Velocity = tpcComponent.Velocity.Add(tpcComponent.ZipVelocity)

	transformComponent.Position = transformComponent.Position.Add(tpcComponent.Velocity.Mul(delta.Seconds()))

	// safeguard falling off the map
	if transformComponent.Position[1] < -1000 {
		transformComponent.Position[1] = 25
	}

	if movementDir.LenSqr() > 0 {
		transformComponent.Orientation = libutils.QuatLookAt(mgl64.Vec3{0, 0, 0}, tpcComponent.ControllerVelocity.Normalize(), mgl64.Vec3{0, 1, 0})
	} else {
		tpcComponent.MovementSpeed = 0
	}
}

func ResolveControllerCollision(entity entities.Entity) {
	cc := entity.GetComponentContainer()
	colliderComponent := cc.ColliderComponent
	transformComponent := cc.TransformComponent
	tpcComponent := cc.ThirdPersonControllerComponent

	if colliderComponent.CollisionInstances != nil {
		contactManifolds := colliderComponent.CollisionInstances[0].ContactManifolds

		separatingVector := combineSeparatingVectors(contactManifolds)
		// separatingVector := minSeparatingVector(contactManifolds)
		transformComponent.Position = transformComponent.Position.Add(separatingVector)
		tpcComponent.Grounded = true
		tpcComponent.Velocity[1] = 0
		tpcComponent.BaseVelocity[1] = 0
		tpcComponent.ZipVelocity = mgl64.Vec3{}
	} else {
		// no collisions were detected (i.e. the ground)
		tpcComponent.Grounded = false
	}
}

func minSeparatingVector(contactManifolds []*collision.ContactManifold) mgl64.Vec3 {
	minVector := contactManifolds[0].Contacts[0].SeparatingVector
	minDistance := contactManifolds[0].Contacts[0].SeparatingDistance

	// one manifold for each object that's being collided with
	for _, contactManifold := range contactManifolds {
		for _, contact := range contactManifold.Contacts {
			if contact.SeparatingDistance < minDistance {
				minVector = contact.SeparatingVector
				minDistance = contact.SeparatingDistance
			}
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
func combineSeparatingVectors(contactManifolds []*collision.ContactManifold) mgl64.Vec3 {
	seenNormals := map[string]mgl64.Vec3{}
	for _, contactManifold := range contactManifolds {
		for _, contact := range contactManifold.Contacts {
			normalHash := hashVec(contact.Normal)
			if v, ok := seenNormals[normalHash]; ok {
				if contact.SeparatingVector.LenSqr() > v.LenSqr() {
					seenNormals[normalHash] = contact.SeparatingVector
				}
				continue
			}

			seenNormals[normalHash] = contact.SeparatingVector
		}
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
func calculateMovementDir(camera entities.Entity, controlVector mgl64.Vec3) mgl64.Vec3 {
	cc := camera.GetComponentContainer()
	forwardVector := cc.TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
	forwardVector = forwardVector.Normalize().Mul(controlVector.Z())
	forwardVector[1] = 0

	rightVector := cc.TransformComponent.Orientation.Rotate(mgl64.Vec3{1, 0, 0})
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
