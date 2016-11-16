package geometry

import "github.com/kkevinchou/ant/lib/math/vector"

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

// TODO: Should I return the internal reference to the points? Or
// return copies? Concern is that points could be manipulated externally
// unintentionally thus breaking the polygon -- extremely hard to debug :(
func (p *Polygon) Points() []Point {
	return p.points
}

// TODO: Might be worth considering caching or constructing the edges at construction time
// as opposed to reconstructing edges each time.  My concern was that edges could be modified
// externally but it may not be a big issue *shrugs*
func (p *Polygon) Edges() []Edge {
	n := len(p.points)
	edges := make([]Edge, n)
	for i, point := range p.points {
		edges[i] = Edge{point, p.points[((i + 1) % n)]}
	}

	return edges
}

// We consider the borders to be inclusive, may be subject to change in the future
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
