package collider

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Capsule struct {
	Radius float64
	Top    mgl64.Vec3
	Bottom mgl64.Vec3
}

func NewCapsule(top, bottom mgl64.Vec3, radius float64) Capsule {
	return Capsule{
		Radius: radius,
		Top:    top,
		Bottom: bottom,
	}
}

func (c Capsule) Transform(position mgl64.Vec3) Capsule {
	newTop := position.Add(c.Top)
	newBottom := position.Add(c.Bottom)
	return NewCapsule(newTop, newBottom, c.Radius)
}

// CreateCapsuleFromMesh creates a capsule centered at the model space origin
func NewCapsuleFromMeshVertices(vertices []mgl64.Vec3) Capsule {
	var minX, minY, minZ, maxX, maxY, maxZ float64

	minX = vertices[0].X()
	maxX = vertices[0].X()
	minY = vertices[0].Y()
	maxY = vertices[0].Y()
	minZ = vertices[0].Z()
	maxZ = vertices[0].Z()
	for _, vertex := range vertices {
		if vertex.X() < minX {
			minX = vertex.X()
		}
		if vertex.X() > maxX {
			maxX = vertex.X()
		}
		if vertex.Y() < minY {
			minY = vertex.Y()
		}
		if vertex.Y() > maxY {
			maxY = vertex.Y()
		}
		if vertex.Z() < minZ {
			minZ = vertex.Z()
		}
		if vertex.Z() > maxZ {
			maxZ = vertex.Z()
		}
	}

	radius := math.Abs(minX)
	radius = math.Max(radius, math.Abs(maxX))
	radius = math.Max(radius, math.Abs(minZ))
	radius = math.Max(radius, math.Abs(maxZ))

	// try our best to construct a capsule that sits above ground
	topYValue := math.Max((2*radius)+1, maxY)

	return Capsule{
		Radius: radius,
		Bottom: mgl64.Vec3{0, radius, 0},
		Top:    mgl64.Vec3{0, topYValue - radius, 0},
	}
}