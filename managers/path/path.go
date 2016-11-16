package path

import (
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/pathing"
)

func sqWithOffset(size, xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{xOffset * size, yOffset * size},
		geometry.Point{xOffset * size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset * size},
	}
	return geometry.NewPolygon(points)
}

func funkyShape1() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{180, 360},
		geometry.Point{180, 420},
		geometry.Point{600, 560},
		geometry.Point{400, 120},
	}
	return geometry.NewPolygon(points)
}

func funkyShape2() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{500, 50},
		geometry.Point{300, 100},
		geometry.Point{400, 100},
	}
	return geometry.NewPolygon(points)
}

func setupNavMesh() *pathing.NavMesh {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 1, 0),
		sqWithOffset(60, 2, 0),
		sqWithOffset(60, 2, 1),
		sqWithOffset(60, 2, 2),
		sqWithOffset(60, 1, 2),
		sqWithOffset(60, 0, 2),
		sqWithOffset(60, 0, 3),
		sqWithOffset(60, 0, 4),
		sqWithOffset(60, 1, 4),
		sqWithOffset(60, 2, 4),
		sqWithOffset(60, 2, 5),
		sqWithOffset(60, 2, 6),
		sqWithOffset(60, 1, 6),
		sqWithOffset(60, 0, 6),
		funkyShape1(),
		funkyShape2(),
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
