package pathing

import (
	"fmt"

	"github.com/kkevinchou/ant/lib/geometry"
)

// Portals are edges that have a deterministic ordering of nodes
// for predictable map lookups.
// TODO: maybe define another node type that doesn't have a polygon pointer?
// chance of bugs occuring where two ordered edges don't equal due to the
// polygon pointer
type Portal struct {
	Node1 Node
	Node2 Node
}

func (p Portal) String() string {
	return fmt.Sprintf("P{%s, %s}", p.Node1, p.Node2)
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
	*RenderComponent
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

	navmesh.RenderComponent = &RenderComponent{navmesh.polygons}

	return navmesh
}
func (nm *NavMesh) Polygons() []*Polygon {
	return nm.polygons
}

func (nm *NavMesh) AddPolygon(geoPolygon *geometry.Polygon) {
	polygon := &Polygon{Polygon: geoPolygon}
	polygon.neighbors = map[*Polygon]Portal{}

	nm.polygons = append(nm.polygons, polygon)
	for _, point1 := range polygon.Points() {
		node1 := Node{X: point1.X, Y: point1.Y, Z: point1.Z, Polygon: polygon}
		for _, point2 := range polygon.Points() {
			node2 := Node{X: point2.X, Y: point2.Y, Z: point2.Z, Polygon: polygon}
			if node1 == node2 {
				continue
			}

			nm.addEdge(node1, node2)

			orderedEdge := CreatePortal(node1, node2)
			// Setting up and detecting nm.portals could probably be done more efficiently
			if len(nm.portals[orderedEdge]) == 1 {
				if nm.portals[orderedEdge][0] != polygon {
					// Set up the neighbors
					polygon.neighbors[nm.portals[orderedEdge][0]] = orderedEdge
					nm.portals[orderedEdge][0].neighbors[polygon] = orderedEdge

					// We've found a portal that connects two pathPolygons, attach the nodes as neighbors
					nm.addEdge(
						Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Z: orderedEdge.Node1.Z, Polygon: nm.portals[orderedEdge][0]},
						Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Z: orderedEdge.Node1.Z, Polygon: polygon},
					)

					nm.addEdge(
						Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Z: orderedEdge.Node1.Z, Polygon: polygon},
						Node{X: orderedEdge.Node1.X, Y: orderedEdge.Node1.Y, Z: orderedEdge.Node1.Z, Polygon: nm.portals[orderedEdge][0]},
					)

					nm.addEdge(
						Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Z: orderedEdge.Node2.Z, Polygon: nm.portals[orderedEdge][0]},
						Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Z: orderedEdge.Node2.Z, Polygon: polygon},
					)

					nm.addEdge(
						Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Z: orderedEdge.Node2.Z, Polygon: polygon},
						Node{X: orderedEdge.Node2.X, Y: orderedEdge.Node2.Y, Z: orderedEdge.Node2.Z, Polygon: nm.portals[orderedEdge][0]},
					)
				}
			} else {
				nm.portals[orderedEdge] = append(nm.portals[orderedEdge], polygon)
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
		v1 := from.Vector3()
		v2 := to.Vector3()

		var length float64
		// TODO: do i need to do this equality check? seems like .Length
		// will already return 0
		if (v1.X == v2.X) && (v1.Y == v2.Y) && (v1.Z == v2.Z) {
			length = 0
		} else {
			length = v1.Sub(v2).Length()
		}

		nm.costs[from][to] = length
		nm.neighbors[from] = append(nm.neighbors[from], to)
	}
}

// TODO: is this function needed?
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
	v1 := from.Vector3()
	v2 := to.Vector3()

	return v1.Sub(v2).Length()
}
