package collision

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/collision/checks"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type Contact struct {
	Point mgl64.Vec3
	// Normal             mgl64.Vec3
	SeparatingVector   mgl64.Vec3
	SeparatingDistance float64
}

type ContactManifold struct {
	Contacts []Contact
}

func CheckCollisionCapsuleTriMesh(capsule collider.Capsule, triangulatedMesh collider.TriMesh) *ContactManifold {
	for _, tri := range triangulatedMesh.Triangles {
		manifold := CheckCollisionCapsuleTriangle(capsule, tri)
		// TODO: handle multiple collided triangles
		if manifold != nil {
			return manifold
		}
	}

	return nil
}

func CheckCollisionCapsuleTriangle(capsule collider.Capsule, triangle collider.Triangle) *ContactManifold {
	closestPoints, closestPointsDistance := checks.ClosestPointsLineVSTriangle(
		collider.Line{P1: capsule.Top, P2: capsule.Bottom},
		triangle,
	)
	closestPointOnTriangle := closestPoints[1]

	if closestPointsDistance < capsule.Radius {
		separatingDistance := capsule.Radius - closestPointsDistance
		return &ContactManifold{
			Contacts: []Contact{
				{
					Point: closestPointOnTriangle,
					// Normal:             triangle.Normal,
					SeparatingVector:   closestPoints[0].Sub(closestPoints[1]).Normalize().Mul(separatingDistance),
					SeparatingDistance: separatingDistance,
				},
			},
		}
	}

	return nil
}

func CheckCollisionSpherePoint(sphere collider.Sphere, point mgl64.Vec3) *ContactManifold {
	lenSq := sphere.Center.Sub(mgl64.Vec3(point)).LenSqr()
	if lenSq < sphere.RadiusSquared {
		return &ContactManifold{
			Contacts: []Contact{
				{
					Point: mgl64.Vec3{point[0], point[1], point[2]},
					// Normal: sphere.Center.Sub(mgl64.Vec3(point)),
				},
			},
		}
	}

	return nil
}
