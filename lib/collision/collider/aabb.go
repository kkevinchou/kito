package collider

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type BoundingBox struct {
	MinVertex mgl64.Vec3
	MaxVertex mgl64.Vec3
}

func (c BoundingBox) Transform(position mgl64.Vec3) BoundingBox {
	return BoundingBox{MinVertex: c.MinVertex.Add(position), MaxVertex: c.MaxVertex.Add(position)}
}

func BoundingBoxFromVertices(vertices []mgl64.Vec3) BoundingBox {
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

	minVertex := mgl64.Vec3{minX, minY, minZ}
	maxVertex := mgl64.Vec3{maxX, maxY, maxZ}

	return BoundingBox{MinVertex: minVertex, MaxVertex: maxVertex}
}

func BoundingBoxFromCapsule(capsule Capsule) BoundingBox {
	diagonal := math.Sqrt(2 * capsule.Radius * capsule.Radius)
	return BoundingBox{MinVertex: capsule.Bottom.Sub(mgl64.Vec3{-diagonal, 0, -diagonal}), MaxVertex: capsule.Bottom.Add(mgl64.Vec3{diagonal, 0, diagonal})}
}
