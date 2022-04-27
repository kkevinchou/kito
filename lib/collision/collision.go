package collision

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/collision/checks"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type Contact struct {
	Point              mgl64.Vec3
	Normal             mgl64.Vec3
	SeparatingVector   mgl64.Vec3
	SeparatingDistance float64
}

type ContactManifold struct {
	TriIndex int
	Contacts []Contact
}

func CheckCollisionCapsuleTriMesh(capsule collider.Capsule, triangulatedMesh collider.TriMesh) []*ContactManifold {
	var contactManifolds []*ContactManifold
	for i, tri := range triangulatedMesh.Triangles {
		contactManifold := CheckCollisionCapsuleTriangle(capsule, tri)
		// TODO: handle multiple collided triangles
		if contactManifold != nil {
			contactManifold.TriIndex = i
			contactManifolds = append(contactManifolds, contactManifold)
		}
	}

	return contactManifolds
}

func CheckCollisionCapsuleTriangle(capsule collider.Capsule, triangle collider.Triangle) *ContactManifold {
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
			// fmt.Println("Wat", time.Now())
			// hacky handling of separating vector pushing the capsule opposite to the triangle normal
			separatingVec = separatingVec.Add(triangle.Normal.Mul(capsule.Radius * 2))
		}
		return &ContactManifold{
			Contacts: []Contact{
				{
					Point:              closestPointTriangle,
					Normal:             triangle.Normal,
					SeparatingVector:   separatingVec,
					SeparatingDistance: separatingDistance,
				},
			},
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
