package primitives

import "github.com/go-gl/mathgl/mgl64"

type Mesh interface {
	Vertices() []mgl64.Vec3
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

type TriMesh struct {
	Triangles []Triangle
}

func NewBoxTriMesh(w, l, h float64) TriMesh {
	halfW := w / 2
	halfL := l / 2
	triMesh := TriMesh{}
	// font
	triMesh.Triangles = append(triMesh.Triangles, NewTriangle(
		[]mgl64.Vec3{{-halfW, 0, halfL}, {halfW, 0, halfL}, {halfW, h, halfL}},
	))
	triMesh.Triangles = append(triMesh.Triangles, NewTriangle(
		[]mgl64.Vec3{{halfW, h, halfL}, {-halfW, h, halfL}, {-halfW, 0, halfL}},
	))
	// back
	// left
	// right
	// bottom
	// top
	triMesh.Triangles = append(triMesh.Triangles, NewTriangle(
		[]mgl64.Vec3{{-halfW, h, halfL}, {halfW, h, halfL}, {halfW, h, -halfL}},
	))
	triMesh.Triangles = append(triMesh.Triangles, NewTriangle(
		[]mgl64.Vec3{{halfW, h, -halfL}, {-halfW, h, -halfL}, {-halfW, h, halfL}},
	))
	return triMesh
}
