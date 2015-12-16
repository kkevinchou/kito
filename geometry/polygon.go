package geometry

import "github.com/kkevinchou/ant/math/vector"

// Assumptions:
// Counter clock-wise winding order of vertices
// Polygons are convex (in the future we will ensure this by splitting nonconvex polygons)

type Point vector.Vector

func (p Point) Vector() vector.Vector {
	return vector.Vector{p.X, p.Y}
}

type Edge struct {
	A Point
	B Point
}

type Polygon struct {
	points []Point
}

func (p *Polygon) Points() []Point {
	return p.points
}

func (p *Polygon) Edges() []Edge {
	n := len(p.points)
	edges := make([]Edge, n)
	for i, point := range p.points {
		edges[i] = Edge{point, p.points[((i + 1) % n)]}
	}

	return edges
}

func (p *Polygon) ContainsPoint(point Point) bool {
	n := len(p.points)

	for i, polygonPoint := range p.points {
		nextPoint := p.points[((i + 1) % n)]
		vector := polygonPoint.Vector()

		affineSegment := nextPoint.Vector().Sub(vector)
		affinePoint := point.Vector().Sub(vector)

		if affineSegment.Cross(affinePoint) > 0 {
			return false
		}
	}
	return true
}

func NewPolygon(p []Point) *Polygon {
	return &Polygon{p}
}
