package pathing

import (
	"testing"

	"github.com/kkevinchou/ant/geometry"
)

func tri1() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{11, 4},
		geometry.Point{13, 10},
		geometry.Point{17, 8},
	}
	return geometry.NewPolygon(points)
}

func tri2() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{13, 10},
		geometry.Point{12, 13},
		geometry.Point{17, 8},
	}
	return geometry.NewPolygon(points)
}

func tri3() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{17, 8},
		geometry.Point{12, 13},
		geometry.Point{21, 7},
	}
	return geometry.NewPolygon(points)
}

func tri4() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{17, 2},
		geometry.Point{17, 8},
		geometry.Point{21, 7},
	}
	return geometry.NewPolygon(points)
}

func TestWithNewApex(t *testing.T) {
	polygons := []*geometry.Polygon{
		tri1(),
		tri2(),
		tri3(),
		tri4(),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{13, 7}, geometry.Point{18, 5})
	expectedPath := []Node{Node{X: 13, Y: 7}, Node{X: 17, Y: 8}, Node{X: 18, Y: 5}}
	assertPathEq(t, expectedPath, path)
}

func TestSmoothing(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithXOffset(0),
		sqWithXOffset(6),
		sqWithXOffset(12),
		sqWithXOffset(18),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{1, 1}, geometry.Point{17, 5})
	expectedPath := []Node{Node{X: 1, Y: 1}, Node{X: 17, Y: 5}}
	assertPathEq(t, expectedPath, path)
}

func TestTwoApexes(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(30, 0, 0),
		sqWithOffset(30, 1, 0),
		sqWithOffset(30, 2, 0),
		sqWithOffset(30, 2, 1),
		sqWithOffset(30, 2, 2),
		sqWithOffset(30, 3, 2),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 0}, geometry.Point{110, 69})
	expectedPath := []Node{Node{X: 0, Y: 0}, Node{X: 60, Y: 30}, Node{X: 90, Y: 60}, Node{X: 110, Y: 69}}
	assertPathEq(t, expectedPath, path)
}

func TestStartNodeOverlapsNode(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(30, 0, 0),
		sqWithOffset(30, 1, 0),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 0}, geometry.Point{50, 20})
	expectedPath := []Node{Node{X: 0, Y: 0}, Node{X: 50, Y: 20}}
	assertPathEq(t, expectedPath, path)
}

func TestGoalNodeOverlapsNode(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(30, 0, 0),
		sqWithOffset(30, 1, 0),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{1, 1}, geometry.Point{30, 30})
	expectedPath := []Node{Node{X: 1, Y: 1}, Node{X: 30, Y: 30}}
	assertPathEq(t, expectedPath, path)
}

func TestStartAndGoalNodeOverlapsNode(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(30, 0, 0),
		sqWithOffset(30, 1, 0),
		sqWithOffset(30, 1, 1),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 0}, geometry.Point{30, 60})
	expectedPath := []Node{Node{X: 0, Y: 0}, Node{X: 30, Y: 30}, Node{X: 30, Y: 60}}
	assertPathEq(t, expectedPath, path)
}

func TestReverseC(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 1, 0),
		sqWithOffset(60, 1, 1),
		sqWithOffset(60, 1, 2),
		sqWithOffset(60, 0, 2),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 0}, geometry.Point{20, 140})
	expectedPath := []Node{Node{X: 0, Y: 0}, Node{X: 60, Y: 60}, Node{X: 60, Y: 120}, Node{X: 20, Y: 140}}
	assertPathEq(t, expectedPath, path)
}

func TestC(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 1, 0),
		sqWithOffset(60, 0, 1),
		sqWithOffset(60, 0, 2),
		sqWithOffset(60, 1, 2),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{80, 20}, geometry.Point{80, 140})
	expectedPath := []Node{Node{X: 80, Y: 20}, Node{X: 60, Y: 60}, Node{X: 60, Y: 120}, Node{X: 80, Y: 140}}
	assertPathEq(t, expectedPath, path)
}

func TestOnEdgeToApex(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 0, 1),
		sqWithOffset(60, -1, 1),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 30}, geometry.Point{-20, 60})
	expectedPath := []Node{Node{X: 0, Y: 30}, Node{X: 0, Y: 60}, Node{X: -20, Y: 60}}
	assertPathEq(t, expectedPath, path)
}

func TestPathDoesNotExist(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 0, 1),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 30}, geometry.Point{61, 0})
	assertPathEq(t, nil, path)
}

func TestStartEqualsGoal(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{0, 30}, geometry.Point{0, 30})
	expectedPath := []Node{Node{X: 0, Y: 30}}
	assertPathEq(t, expectedPath, path)
}

func sqWithOffset(size, xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{xOffset * size, yOffset * size},
		geometry.Point{xOffset * size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset * size},
	}
	return geometry.NewPolygon(points)
}

func sqWithXOffset(offset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{offset + 0, 0},
		geometry.Point{offset + 0, 6},
		geometry.Point{offset + 6, 6},
		geometry.Point{offset + 6, 0},
	}
	return geometry.NewPolygon(points)
}

func assertPathEq(t *testing.T, expected, actual []Node) {
	if len(actual) != len(expected) {
		t.Fatalf("Expected: %s Actual: %s", expected, actual)
	}

	for i, node := range actual {
		if node != expected[i] {
			t.Fatalf("Expected: %s Actual: %s", expected, actual)
		}
	}

}
