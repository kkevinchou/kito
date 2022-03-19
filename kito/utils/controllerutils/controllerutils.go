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
	equalThreshold float64 = 1e-5
)

var (
	accelerationDueToGravity = mgl64.Vec3{0, -gravity, 0}
)

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

	// handle controller movement
	movementDir := calculateMovementDir(camera, controlVector)
	tpcComponent.MovementSpeed = computeMoveSpeed(tpcComponent.MovementSpeed)
	tpcComponent.ControllerVelocity = movementDir.Mul(tpcComponent.MovementSpeed)

	tpcComponent.BaseVelocity = tpcComponent.BaseVelocity.Add(accelerationDueToGravity.Mul(delta.Seconds()))
	tpcComponent.Velocity = tpcComponent.BaseVelocity.Add(tpcComponent.ControllerVelocity)
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

func combineSeparatingVectors(contactManifolds []*collision.ContactManifold) mgl64.Vec3 {
	// only add separating vectors which we haven't seen before. ideally
	// this should handle cases where separating vectors are a basis of another
	// and avoid "overcounting" separating vectors
	seenSeparatingVectors := []mgl64.Vec3{}
	triIndices := []int{}
	var separatingVector mgl64.Vec3
	for _, contactManifold := range contactManifolds {
		for _, contact := range contactManifold.Contacts {
			seen := false
			for _, v := range seenSeparatingVectors {
				// if contact.SeparatingVector.ApproxEqual(v) {
				if contact.SeparatingVector.ApproxEqualThreshold(v, equalThreshold) {
					seen = true
					break
				}
			}
			if !seen {
				triIndices = append(triIndices, contactManifold.TriIndex)
				separatingVector = separatingVector.Add(contact.SeparatingVector)
				seenSeparatingVectors = append(seenSeparatingVectors, contact.SeparatingVector)
			}
		}
	}

	if len(seenSeparatingVectors) > 2 {
		fmt.Println("-------------------------")
		fmt.Println("seenSeparatingVectors len > 2")
		fmt.Println("triIndices")
		for _, t := range triIndices {
			fmt.Println(t)
		}
		fmt.Println("separating vecs")
		for _, v := range seenSeparatingVectors {
			fmt.Println(v)
		}
	}
	return separatingVector
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
	cameraComponentContainer := camera.GetComponentContainer()
	forwardVector := cameraComponentContainer.TransformComponent.Orientation.Rotate(mgl64.Vec3{0, 0, -1})
	forwardVector[1] = 0
	forwardVector = forwardVector.Normalize().Mul(controlVector.Z())

	rightVector := cameraComponentContainer.TransformComponent.Orientation.Rotate(mgl64.Vec3{1, 0, 0})
	rightVector[1] = 0
	rightVector = rightVector.Normalize().Mul(controlVector.X())

	movementDir := forwardVector.Add(rightVector)
	if movementDir.LenSqr() > 0 {
		return movementDir.Normalize()
	}

	return movementDir
}
