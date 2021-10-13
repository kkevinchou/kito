package primitives

import "github.com/go-gl/mathgl/mgl64"

type TriangulatedMesh struct {
	Triangles []Triangle
}

type Triangle struct {
	Normal mgl64.Vec3
	Points []mgl64.Vec3
}

func NewTriangle(points []mgl64.Vec3) Triangle {
	seg1 := points[1].Sub(points[0])
	seg2 := points[2].Sub(points[0])
	normal := seg1.Cross(seg2).Normalize()
	return Triangle{
		Points: points,
		Normal: normal,
	}
}

type Line struct {
	P1 mgl64.Vec3
	P2 mgl64.Vec3
}
