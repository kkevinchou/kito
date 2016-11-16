package pathing

import (
	"fmt"

	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
)

// Assumptions for pathfinding
// 1. Polygons are non overlapping - though they can share vertices
// 2. Polygons are convex

type Polygon struct {
	*geometry.Polygon
	neighbors map[*Polygon]Portal
}

type Node struct {
	X       float64
	Y       float64
	Polygon *Polygon
}

func (n Node) PositionEquals(other Node) bool {
	return (n.X == other.X) && (n.Y == other.Y)
}

func (n Node) Vector() vector.Vector {
	return vector.Vector{n.X, n.Y}
}

func (n Node) Less(other Node) bool {
	if n.X < other.X {
		return true
	} else if n.X > other.X {
		return false
	}

	if n.Y < other.Y {
		return true
	} else if n.Y > other.Y {
		return false
	}

	return true
}

func (n Node) String() string {
	return fmt.Sprintf("N[%.1f, %.1f]", n.X, n.Y)
}
