package collision

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision/checks"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type Contact struct {
	EntityID       *int
	SourceEntityID *int

	TriIndex           *int
	Point              mgl64.Vec3
	Normal             mgl64.Vec3
	SeparatingVector   mgl64.Vec3
	SeparatingDistance float64
}

func (c *Contact) String() string {
	var result = fmt.Sprintf("{ EntityID: %d, TriIndex: %d ", *c.EntityID, *c.TriIndex)
	result += fmt.Sprintf("[ SV: %s, N: %s, D: %.3f D2: %.3f]", utils.PPrintVec(c.SeparatingVector), utils.PPrintVec(c.Normal), c.SeparatingDistance, c.SeparatingVector.Len())
	result += " }"
	return result
}

func CheckCollisionCapsuleTriMesh(capsule collider.Capsule, triangulatedMesh collider.TriMesh) []*Contact {
	var contacts []*Contact
	for i, tri := range triangulatedMesh.Triangles {
		if triContact := CheckCollisionCapsuleTriangle(capsule, tri); triContact != nil {
			index := i
			triContact.TriIndex = &index
			contacts = append(contacts, triContact)
		}
	}

	return contacts
}

func CheckCollisionCapsuleTriangle(capsule collider.Capsule, triangle collider.Triangle) *Contact {
	closestPoints, closestPointsDistance := checks.ClosestPointsLineVSTriangle(
		collider.Line{P1: capsule.Top, P2: capsule.Bottom},
		triangle,
	)
	// closestPointCapsule := closestPoints[0]
	closestPointTriangle := closestPoints[1]

	if closestPointsDistance < capsule.Radius {
		separatingDistance := capsule.Radius - closestPointsDistance
		separatingVec := closestPoints[0].Sub(closestPoints[1]).Normalize().Mul(separatingDistance)
		if separatingVec.Dot(triangle.Normal) < 0 {
			// TODO(kevin): not sure if this is right, might want to revisit
			// hacky handling of separating vector pushing the capsule opposite to the triangle normal
			separatingVec = separatingVec.Add(triangle.Normal.Mul(capsule.Radius * 2))
			separatingDistance = separatingVec.Len()
		}
		return &Contact{
			Point:              closestPointTriangle,
			Normal:             triangle.Normal,
			SeparatingVector:   separatingVec,
			SeparatingDistance: separatingDistance,
		}
	}

	return nil
}

// func CheckCollisionSpherePoint(sphere collider.Sphere, point mgl64.Vec3) *ContactManifold {
// 	lenSq := sphere.Center.Sub(mgl64.Vec3(point)).LenSqr()
// 	if lenSq < sphere.RadiusSquared {
// 		return &ContactManifold{
// 			Contacts: []Contact{
// 				{
// 					Point: mgl64.Vec3{point[0], point[1], point[2]},
// 					// Normal: sphere.Center.Sub(mgl64.Vec3(point)),
// 				},
// 			},
// 		}
// 	}

// 	return nil
// }
