package path

import (
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/pathing"
)

func sqWithOffset(size, xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: xOffset * size, Y: 0, Z: yOffset * size},
		geometry.Point{X: xOffset * size, Y: 0, Z: yOffset*size + size},
		geometry.Point{X: xOffset*size + size, Y: 0, Z: yOffset*size + size},
		geometry.Point{X: xOffset*size + size, Y: 0, Z: yOffset * size},
	}
	return geometry.NewPolygon(points)
}

func funkyShape1() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 180, Y: 0, Z: 360},
		geometry.Point{X: 180, Y: 0, Z: 420},
		geometry.Point{X: 600, Y: 0, Z: 560},
		geometry.Point{X: 400, Y: 0, Z: 120},
	}
	return geometry.NewPolygon(points)
}

func funkyShape2() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{X: 500, Y: 0, Z: 50},
		geometry.Point{X: 300, Y: 0, Z: 100},
		geometry.Point{X: 400, Y: 0, Z: 100},
	}
	return geometry.NewPolygon(points)
}

func setupNavMesh() *pathing.NavMesh {
	// polygons := []*geometry.Polygon{
	// 	sqWithOffset(60, 0, 0),
	// 	sqWithOffset(60, 1, 0),
	// 	sqWithOffset(60, 2, 0),
	// 	sqWithOffset(60, 2, 1),
	// 	sqWithOffset(60, 2, 2),
	// 	sqWithOffset(60, 1, 2),
	// 	sqWithOffset(60, 0, 2),
	// 	sqWithOffset(60, 0, 3),
	// 	sqWithOffset(60, 0, 4),
	// 	sqWithOffset(60, 1, 4),
	// 	sqWithOffset(60, 2, 4),
	// 	sqWithOffset(60, 2, 5),
	// 	sqWithOffset(60, 2, 6),
	// 	sqWithOffset(60, 1, 6),
	// 	sqWithOffset(60, 0, 6),
	// 	funkyShape1(),
	// 	funkyShape2(),
	// }

	points := []geometry.Point{
		geometry.Point{X: -10, Y: 0, Z: -10},
		geometry.Point{X: -10, Y: 0, Z: 10},
		geometry.Point{X: 10, Y: 0, Z: 10},
		geometry.Point{X: 10, Y: 0, Z: -10},
	}

	polygons := []*geometry.Polygon{
		geometry.NewPolygon(points),
	}

	return pathing.ConstructNavMesh(polygons)
}

type Manager struct {
	planner pathing.Planner
	navMesh *pathing.NavMesh
}

func (m *Manager) FindPath(start, goal geometry.Point) []pathing.Node {
	return m.planner.FindPath(start, goal)
}

func (m *Manager) NavMesh() *pathing.NavMesh {
	return m.navMesh
}

func NewManager() *Manager {
	p := pathing.Planner{}
	navMesh := setupNavMesh()
	p.SetNavMesh(navMesh)
	return &Manager{planner: p, navMesh: navMesh}
}
