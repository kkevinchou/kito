package pathing

import (
	"fmt"
	"testing"

	"github.com/kkevinchou/ant/geometry"
	"github.com/kkevinchou/ant/pathing"
)

//
//
//
//
//
//
//
//
//
//
//
//
//
//

func createPolygon(xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{xOffset + 0, yOffset + 0},
		geometry.Point{xOffset + 0, yOffset + 6},
		geometry.Point{xOffset + 6, yOffset + 6},
		geometry.Point{xOffset + 6, yOffset + 0},
	}
	return geometry.NewPolygon(points)
}

func TestPriorityQueue(t *testing.T) {
	polygons := []*geometry.Polygon{
		createPolygon(0, 0),
		createPolygon(6, 0),
		createPolygon(12, 0),
	}

	navmesh := pathing.ConstructNavMesh(polygons)
	p := pathing.Planner{}
	p.SetNavMesh(navmesh)

	fmt.Println(p.FindPath(geometry.Point{1, 5}, geometry.Point{17, 5}))
	// fmt.Println(graph.Neighbors(pathing.Node{6, 6}))
	// n1 := CreateNode(0, 0)
	// n2 := CreateNode(1, 1)
	// p := CreatePlanner()
	t.Fail()
}
