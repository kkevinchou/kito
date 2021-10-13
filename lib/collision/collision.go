package collision

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/collision/checks"
	"github.com/kkevinchou/kito/lib/collision/primitives"
)

type Contact struct {
	Point              mgl64.Vec3
	Normal             mgl64.Vec3
	SeparatingDistance float64
}

type ContactManifold struct {
	Contacts []Contact
}

func CheckCollision(capsule primitives.Capsule, triangulatedMesh primitives.TriangulatedMesh) *ContactManifold {
	for _, tri := range triangulatedMesh.Triangles {
		manifold := CheckCollisionCapsuleTriangle(capsule, tri)
		// TODO: handle multiple collided triangles
		if manifold != nil {
			return manifold
		}
	}

	return nil
}

func CheckCollisionCapsuleTriangle(capsule primitives.Capsule, triangle primitives.Triangle) *ContactManifold {
	closestPoints, closestPointsDistance := checks.ClosestPointsLineVSTriangle(
		primitives.Line{P1: capsule.Top, P2: capsule.Bottom},
		triangle,
	)
	closestPointOnTriangle := closestPoints[1]

	if closestPointsDistance < capsule.Radius {
		return &ContactManifold{
			Contacts: []Contact{
				{
					Point:              closestPointOnTriangle,
					Normal:             triangle.Normal,
					SeparatingDistance: closestPointsDistance,
				},
			},
		}
	}

	return nil
}

func CheckCollisionSpherePoint(sphere primitives.Sphere, point mgl64.Vec3) *ContactManifold {
	lenSq := sphere.Center.Sub(mgl64.Vec3(point)).LenSqr()
	if lenSq < sphere.RadiusSquared {
		return &ContactManifold{
			Contacts: []Contact{
				{
					Point:  mgl64.Vec3{point[0], point[1], point[2]},
					Normal: sphere.Center.Sub(mgl64.Vec3(point)),
				},
			},
		}
	}

	return nil
}
