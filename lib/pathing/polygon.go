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

// TODO: just store x,y,z in a vector??
type Node struct {
	X       float64
	Y       float64
	Z       float64
	Polygon *Polygon
}

func (n Node) PositionEquals(other Node) bool {
	return (n.X == other.X) && (n.Y == other.Y) && (n.Z == other.Z)
}

func (n Node) Vector3() vector.Vector3 {
	return vector.Vector3{n.X, n.Y, n.Z}
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

	if n.Z < other.Z {
		return true
	} else if n.Z > other.Z {
		return false
	}

	return true
}

func (n Node) String() string {
	return fmt.Sprintf("N[%.1f, %.1f, %.1f]", n.X, n.Y, n.Z)
}
