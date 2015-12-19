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

type Polygon struct {
	*geometry.Polygon
	neighbors map[*Polygon]Portal
}

type Node struct {
	X       float64
	Y       float64
	Polygon *Polygon
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

// Portals are edges that have a deterministic ordering of nodes
// for predictable map lookups.
// TODO: maybe define another node type that doesn't have a polygon pointer?
// chance of bugs occuring where two ordered edges don't equal due to the
// polygon pointer
type Portal struct {
	Node1 Node
	Node2 Node
}

func CreatePortal(node1, node2 Node) Portal {
	node1.Polygon = nil
	node2.Polygon = nil
	if node1.Less(node2) {
		return Portal{Node1: node1, Node2: node2}
	} else {
		return Portal{Node1: node2, Node2: node1}
	}
}

type NavMesh struct {
	neighbors map[Node][]Node
	costs     map[Node]map[Node]float64
	polygons  []*Polygon
	portals   map[Portal][]*Polygon
}

func (nm *NavMesh) Polygons() []*Polygon {
	return nm.polygons
}

func ConstructNavMesh(polygons []*geometry.Polygon) *NavMesh {
	navmesh := &NavMesh{
		neighbors: map[Node][]Node{},
		costs:     map[Node]map[Node]float64{},
		portals:   map[Portal][]*Polygon{},
	}

	for _, polygon := range polygons {
		navmesh.AddPolygon(polygon)
	}

	return navmesh
}

func (nm *NavMesh) AddPolygon(geoPolygon *geometry.Polygon) {
	polygon := &Polygon{Polygon: geoPolygon}
	polygon.neighbors = map[*Polygon]Portal{}

	nm.polygons = append(nm.polygons, polygon)
	for _, point1 := range polygon.Points() {
		node1 := Node{X: point1.X, Y: point1.Y, Polygon: polygon}
		for _, point2 := range polygon.Points() {
			node2 := Node{X: point2.X, Y: point2.Y, Polygon: polygon}
			if node1 != node2 {
				nm.addEdge(node1, node2)

				orderedEdge := CreatePortal(node1, node2)
				// Setting up and detecting nm.portals could probably be done more efficiently
				if len(nm.portals[orderedEdge]) == 1 {
					if nm.portals[orderedEdge][0] != polygon {
						polygon.neighbors[nm.portals[orderedEdge][0]] = orderedEdge
						nm.portals[orderedEdge][0].neighbors[polygon] = orderedEdge
						// We've found a portal that connects two pathPolygons, attach the nodes as neighbors
						nm.addEdge(
							Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Polygon: nm.portals[orderedEdge][0]},
							Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Polygon: polygon},
						)

						nm.addEdge(
							Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Polygon: polygon},
							Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Polygon: nm.portals[orderedEdge][0]},
						)

						nm.addEdge(
							Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Polygon: nm.portals[orderedEdge][0]},
							Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Polygon: polygon},
						)

						nm.addEdge(
							Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Polygon: polygon},
							Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Polygon: nm.portals[orderedEdge][0]},
						)
					}
				} else {
					nm.portals[orderedEdge] = append(nm.portals[orderedEdge], polygon)
				}
			}
		}
	}
}

func (nm *NavMesh) addEdge(from, to Node) {
	if _, ok := nm.neighbors[from]; !ok {
		nm.neighbors[from] = []Node{}
		nm.costs[from] = map[Node]float64{}
	}

	if _, ok := nm.costs[from][to]; !ok {
		v1 := vector.Vector{X: from.X, Y: from.Y}
		v2 := vector.Vector{X: to.X, Y: to.Y}

		var length float64
		if (v1.X == v2.X) && (v1.Y == v1.Y) {
			length = 0
		} else {
			length = v1.Sub(v2).Length()
		}

		nm.costs[from][to] = length
		nm.neighbors[from] = append(nm.neighbors[from], to)
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

	// Initialize
	frontier := priorityqueue.New()
	cameFrom := map[Node]Node{}
	costSoFar := map[Node]float64{startNode: 0}
	explored := map[Node]bool{}

	// Find which polygon our start node lies in
	for _, polygon := range p.navmesh.Polygons() {
		if polygon.ContainsPoint(start) {
			for _, point := range polygon.Points() {
				node := Node{X: point.X, Y: point.Y, Polygon: polygon}
				cost := point.Vector().Sub(start.Vector()).Length()
				startNode.Polygon = polygon
				cameFrom[node] = startNode
				costSoFar[node] = cost

				// Initialize the frontier with each of the neighbors
				// of the start node within the polygon
				frontier.Push(node, cost)
			}
			break
		}
	}

	if frontier.Empty() {
		// start does not reside in any part of our navmesh
		return nil
	}

	// Find which polygon our goal node lies in
	var goalPolygon *Polygon
	for _, polygon := range p.navmesh.Polygons() {
		if polygon.ContainsPoint(goal) {
			goalNode.Polygon = polygon
			goalPolygon = polygon
			break
		}
	}

	if goalPolygon == nil {
		// goal does not reside in any part of our navmesh
		return nil
	}

	// If we have a direct path from start to goal, return it
	if startNode.Polygon == goalNode.Polygon {
		return []Node{startNode, goalNode}
	}

	// Set the goal node as the neighbor of each node in the goal
	// polygon
	goalNeighbors := map[Node]bool{}
	for _, point := range goalPolygon.Points() {
		node := Node{X: point.X, Y: point.Y, Polygon: goalPolygon}
		goalNeighbors[node] = true
	}

	// Start searching for a path!
	for !frontier.Empty() {
		current := frontier.Pop().(Node)

		if current == goalNode {
			break
		}

		explored[current] = true

		neighbors := p.navmesh.Neighbors(current)
		if _, ok := goalNeighbors[current]; ok {
			neighbors = append(neighbors, goalNode)
		}

		for _, neighbor := range neighbors {
			if _, ok := explored[neighbor]; ok {
				continue
			}

			newCost := costSoFar[current] + p.navmesh.Cost(current, neighbor)
			if cost, ok := costSoFar[neighbor]; !ok || newCost < cost {
				costSoFar[neighbor] = newCost
				frontier.Push(neighbor, newCost+p.navmesh.HeuristicCost(goalNode, neighbor))
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
