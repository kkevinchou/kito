package pathing

import (
	"testing"

	"github.com/kkevinchou/kito/lib/geometry"
)

func tri1() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 11, Y: 0, Z: 4},
		geometry.Point{X: 13, Y: 0, Z: 10},
		geometry.Point{X: 17, Y: 0, Z: 8},
	}
	return geometry.NewPolygon(points)
}

func tri2() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 13, Y: 0, Z: 10},
		geometry.Point{X: 12, Y: 0, Z: 13},
		geometry.Point{X: 17, Y: 0, Z: 8},
	}
	return geometry.NewPolygon(points)
}

func tri3() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 17, Y: 0, Z: 8},
		geometry.Point{X: 12, Y: 0, Z: 13},
		geometry.Point{X: 21, Y: 0, Z: 7},
	}
	return geometry.NewPolygon(points)
}

func tri4() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 17, Y: 0, Z: 2},
		geometry.Point{X: 17, Y: 0, Z: 8},
		geometry.Point{X: 21, Y: 0, Z: 7},
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

	path := p.FindPath(geometry.Point{X: 13, Y: 0, Z: 7}, geometry.Point{X: 18, Y: 0, Z: 5})
	expectedPath := []geometry.Point{geometry.Point{X: 13, Y: 0, Z: 7}, geometry.Point{X: 17, Y: 0, Z: 8}, geometry.Point{X: 18, Y: 0, Z: 5}}
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

	path := p.FindPath(geometry.Point{X: 1, Y: 0, Z: 1}, geometry.Point{X: 17, Y: 0, Z: 5})
	expectedPath := []geometry.Point{geometry.Point{X: 1, Y: 0, Z: 1}, geometry.Point{X: 17, Y: 0, Z: 5}}
	assertPathEq(t, expectedPath, path)
}

// X X X
//     X
//     X X
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 110, Y: 0, Z: 69})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0}, geometry.Point{X: 60, Y: 0, Z: 30}, geometry.Point{X: 90, Y: 0, Z: 60}, geometry.Point{X: 110, Y: 0, Z: 69}}
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 50, Y: 0, Z: 20})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 50, Y: 0, Z: 20}}
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

	path := p.FindPath(geometry.Point{X: 1, Y: 0, Z: 1}, geometry.Point{X: 30, Y: 0, Z: 30})
	expectedPath := []geometry.Point{geometry.Point{X: 1, Y: 0, Z: 1}, geometry.Point{X: 30, Y: 0, Z: 30}}
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 30, Y: 0, Z: 60})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 30, Y: 0, Z: 30}, geometry.Point{X: 30, Y: 0, Z: 60}}
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 0}, geometry.Point{X: 20, Y: 0, Z: 140})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0}, geometry.Point{X: 60, Y: 0, Z: 60}, geometry.Point{X: 60, Y: 0, Z: 120}, geometry.Point{X: 20, Y: 0, Z: 140}}
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

	path := p.FindPath(geometry.Point{X: 80, Y: 0, Z: 20}, geometry.Point{X: 80, Y: 0, Z: 140})
	expectedPath := []geometry.Point{geometry.Point{X: 80, Y: 0, Z: 20}, geometry.Point{X: 60, Y: 0, Z: 60}, geometry.Point{X: 60, Y: 0, Z: 120}, geometry.Point{X: 80, Y: 0, Z: 140}}
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 30}, geometry.Point{X: -20, Y: 0, Z: 60})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0, Z: 30}, geometry.Point{X: 0, Y: 0, Z: 60}, geometry.Point{X: -20, Y: 0, Z: 60}}
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

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 30}, geometry.Point{X: 61, Y: 0, Z: 0})
	assertPathEq(t, nil, path)
}

func TestStartEqualsGoal(t *testing.T) {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
	}

	navmesh := ConstructNavMesh(polygons)
	p := Planner{}
	p.SetNavMesh(navmesh)

	path := p.FindPath(geometry.Point{X: 0, Y: 0, Z: 30}, geometry.Point{X: 0, Y: 0, Z: 30})
	expectedPath := []geometry.Point{geometry.Point{X: 0, Y: 0, Z: 30}}
	assertPathEq(t, expectedPath, path)
}

func sqWithOffset(size, xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: xOffset * size, Y: 0, Z: yOffset * size},
		geometry.Point{X: xOffset * size, Y: 0, Z: yOffset*size + size},
		geometry.Point{X: xOffset*size + size, Y: 0, Z: yOffset*size + size},
		geometry.Point{X: xOffset*size + size, Y: 0, Z: yOffset * size},
	}
	return geometry.NewPolygon(points)
}

func sqWithXOffset(offset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: offset + 0, Y: 0, Z: 0},
		geometry.Point{X: offset + 0, Y: 0, Z: 6},
		geometry.Point{X: offset + 6, Y: 0, Z: 6},
		geometry.Point{X: offset + 6, Y: 0, Z: 0},
	}
	return geometry.NewPolygon(points)
}

func assertPathEq(t *testing.T, expected, actual []geometry.Point) {
	if len(actual) != len(expected) {
		t.Fatalf("Expected: %v Actual: %v", expected, actual)
	}

	for i, point := range actual {
		if point != expected[i] {
			t.Fatalf("Expected: %v Actual: %v", expected, actual)
		}
	}
}
