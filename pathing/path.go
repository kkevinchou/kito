package pathing

import (
	"fmt"

	"github.com/kkevinchou/ant/geometry"
	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/util"
)

// Assumptions for pathfinding
// 1. Polygons are non overlapping - though they can share vertices
// 2. Polygons are convex

type Node vector.Vector

func (n Node) String() string {
	return fmt.Sprintf("N[%.1f, %.1f]", n.X, n.Y)
}

type Edge struct {
	From   Node
	To     Node
	length float64
}

func (e Edge) String() string {
	return fmt.Sprintf("E{%s -> %s | %.1f}", e.From, e.To, e.length)
}

func CreateEdge(from, to Node) Edge {
	v1 := vector.Vector{X: from.X, Y: from.Y}
	v2 := vector.Vector{X: to.X, Y: to.Y}

	length := v1.Sub(v2).Length()

	return Edge{From: from, To: to, length: length}
}

type NavMesh struct {
	neighbors map[Node][]Node
	costs     map[Node]map[Node]float64
	polygons  []*geometry.Polygon
}

func (nm *NavMesh) Polygons() []*geometry.Polygon {
	return nm.polygons
}

func ConstructNavMesh(polygons []*geometry.Polygon) *NavMesh {
	navmesh := &NavMesh{
		neighbors: map[Node][]Node{},
		costs:     map[Node]map[Node]float64{},
	}

	for _, polygon := range polygons {
		navmesh.AddPolygon(polygon)
	}

	return navmesh
}

func (nm *NavMesh) addEdge(from, to Node) {
	if _, ok := nm.neighbors[from]; !ok {
		nm.neighbors[from] = []Node{}
		nm.costs[from] = map[Node]float64{}
	}

	edge := CreateEdge(from, to)

	if _, ok := nm.costs[from][to]; !ok {
		nm.costs[from][to] = edge.length
		nm.neighbors[from] = append(nm.neighbors[from], to)
	}
}

func (nm *NavMesh) AddPolygon(polygon *geometry.Polygon) {
	nm.polygons = append(nm.polygons, polygon)
	for _, point1 := range polygon.Points() {
		node1 := Node{point1.X, point1.Y}
		for _, point2 := range polygon.Points() {
			node2 := Node{point2.X, point2.Y}
			if node1 != node2 {
				nm.addEdge(node1, node2)
			}
		}
	}
}

func (nm *NavMesh) Neighbors(node Node) []Node {
	if neighbors, ok := nm.neighbors[node]; ok {
		return neighbors
	}
	return []Node{}
}

func (nm *NavMesh) Cost(from, to Node) float64 {
	return nm.costs[from][to]
}

func (nm *NavMesh) HeuristicCost(from, to Node) float64 {
	v1 := vector.Vector{X: from.X, Y: from.Y}
	v2 := vector.Vector{X: to.X, Y: to.Y}

	return v1.Sub(v2).Length()
}

type Planner struct {
	navmesh *NavMesh
}

func (p *Planner) SetNavMesh(navmesh *NavMesh) {
	p.navmesh = navmesh
}

func (p *Planner) FindPath(start geometry.Point, goal geometry.Point) []Node {
	startNode, goalNode := Node{X: start.X, Y: start.Y}, Node{X: goal.X, Y: goal.Y}

	frontier := priorityqueue.New()
	cameFrom := map[Node]Node{}
	costSoFar := map[Node]float64{startNode: 0}

	// Find which polygon our start and goal lies in

	for _, polygon := range p.navmesh.Polygons() {
		if polygon.ContainsPoint(start) {
			for _, point := range polygon.Points() {
				node := Node{point.X, point.Y}
				cost := point.Vector().Sub(start.Vector()).Length()
				cameFrom[node] = startNode
				costSoFar[node] = cost

				frontier.Push(node, cost)
			}
			break
		}
	}

	if frontier.Empty() {
		// start does not reside in any part of our navmesh
		return nil
	}

	var goalPolygon *geometry.Polygon
	for _, polygon := range p.navmesh.Polygons() {
		if polygon.ContainsPoint(goal) {
			goalPolygon = polygon
			break
		}
	}

	if goalPolygon == nil {
		// goal does not reside in any part of our navmesh
		return nil
	}

	goalNeighbors := map[Node]bool{}
	for _, point := range goalPolygon.Points() {
		node := Node{X: point.X, Y: point.Y}
		goalNeighbors[node] = true
	}

	for !frontier.Empty() {
		current := frontier.Pop().(Node)

		if current == goalNode {
			break
		}

		neighbors := p.navmesh.Neighbors(current)
		if _, ok := goalNeighbors[current]; ok {
			neighbors = append(neighbors, goalNode)
		}

		for _, neighbor := range neighbors {
			newCost := costSoFar[current] + p.navmesh.Cost(current, neighbor)
			if cost, ok := costSoFar[neighbor]; !ok || newCost < cost {
				costSoFar[neighbor] = newCost
				frontier.Push(neighbor, newCost+p.navmesh.HeuristicCost(current, neighbor))
				cameFrom[neighbor] = current
			}
		}
	}

	if _, ok := cameFrom[goalNode]; !ok {
		// Could not find a path to the goal node
		return nil
	}

	path := []Node{}
	var ok bool
	pathNode := goalNode

	for {
		path = append(path, pathNode)
		if pathNode, ok = cameFrom[pathNode]; !ok {
			break
		}
	}

	reversePath := make([]Node, len(path))
	for i := 0; i < len(path); i++ {
		reversePath[len(path)-1-i] = path[i]
	}

	return reversePath
}
