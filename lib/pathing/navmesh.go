package pathing

import (
	"fmt"

	"github.com/kkevinchou/kito/lib/geometry"
)

type NavNode struct {
	Point   geometry.Point
	Polygon *geometry.Polygon
}

// Portals are edges that have a deterministic ordering of nodes
// for predictable map lookups.
// TODO: maybe define another node type that doesn't have a polygon pointer?
// chance of bugs occuring where two ordered edges don't equal due to the
// polygon pointer
type Portal struct {
	Point1 geometry.Point
	Point2 geometry.Point
}

func (p Portal) String() string {
	return fmt.Sprintf("P{%v, %v}", p.Point1, p.Point2)
}

func CreatePortal(point1, point2 geometry.Point) Portal {
	return Portal{Point1: point1, Point2: point2}
}

type NavMesh struct {
	neighbors        map[NavNode][]NavNode
	costs            map[NavNode]map[NavNode]float64
	polygons         []*geometry.Polygon
	portalToPolygons map[Portal][]*geometry.Polygon

	polyPairToPortal map[*geometry.Polygon]map[*geometry.Polygon]Portal
	pointToPolygons  map[geometry.Point][]*geometry.Polygon

	*RenderComponent
}

func ConstructNavMesh(polygons []*geometry.Polygon) *NavMesh {
	navmesh := &NavMesh{
		neighbors:        map[NavNode][]NavNode{},
		costs:            map[NavNode]map[NavNode]float64{},
		portalToPolygons: map[Portal][]*geometry.Polygon{},
		polyPairToPortal: map[*geometry.Polygon]map[*geometry.Polygon]Portal{},
		pointToPolygons:  map[geometry.Point][]*geometry.Polygon{},
	}

	for _, polygon := range polygons {
		navmesh.AddPolygon(polygon)
	}

	navmesh.RenderComponent = &RenderComponent{
		RenderData: &NavMeshRenderData{
			ID:      "tile",
			Visible: true,
		},
	}

	return navmesh
}
func (nm *NavMesh) Polygons() []*geometry.Polygon {
	return nm.polygons
}

func (nm *NavMesh) AddPolygon(polygon *geometry.Polygon) {
	nm.polygons = append(nm.polygons, polygon)
	for _, point1 := range polygon.Points() {
		nm.pointToPolygons[point1] = append(nm.pointToPolygons[point1], polygon)

		for _, point2 := range polygon.Points() {
			if point1 == point2 {
				continue
			}

			navNode1 := NavNode{Point: point1, Polygon: polygon}
			navNode2 := NavNode{Point: point2, Polygon: polygon}

			// TODO: handle duplicate edges being added
			nm.addEdge(navNode1, navNode2)
			nm.addEdge(navNode2, navNode1)

			portal := Portal{Point1: point1, Point2: point2}

			// TODO: Setting up and detecting nm.portalToPolygons could probably be done more efficiently

			// Found one half of the portal, complete the other half
			if len(nm.portalToPolygons[portal]) == 1 {
				polyWithSharedPortal := nm.portalToPolygons[portal][0]
				if polyWithSharedPortal != polygon {
					// Set up the neighbors
					if _, ok := nm.polyPairToPortal[polygon]; !ok {
						nm.polyPairToPortal[polygon] = map[*geometry.Polygon]Portal{}
					}
					if _, ok := nm.polyPairToPortal[polyWithSharedPortal]; !ok {
						nm.polyPairToPortal[polyWithSharedPortal] = map[*geometry.Polygon]Portal{}
					}

					nm.polyPairToPortal[polygon][polyWithSharedPortal] = portal
					nm.polyPairToPortal[polyWithSharedPortal][polygon] = portal

					// Set points that lie on the same portal to be neighbors to one another

					navNode1 := NavNode{Point: point1, Polygon: polygon}
					otherNavNode1 := NavNode{Point: point1, Polygon: polyWithSharedPortal}

					nm.neighbors[navNode1] = append(
						nm.neighbors[navNode1],
						otherNavNode1,
					)

					nm.neighbors[otherNavNode1] = append(
						nm.neighbors[otherNavNode1],
						navNode1,
					)

					navNode2 := NavNode{Point: point2, Polygon: polygon}
					otherNavNode2 := NavNode{Point: point2, Polygon: polyWithSharedPortal}

					nm.neighbors[navNode2] = append(
						nm.neighbors[navNode2],
						otherNavNode2,
					)

					nm.neighbors[otherNavNode2] = append(
						nm.neighbors[otherNavNode2],
						navNode2,
					)
				}
			}
			nm.portalToPolygons[portal] = append(nm.portalToPolygons[portal], polygon)
		}
	}
}

func (nm *NavMesh) GetPortalFromPolyPair(poly1, poly2 *geometry.Polygon) Portal {
	return Portal{}
}

// TODO: seems weird that we need to update costs here,
// an edge will always have the shortest cost between the two points right?
func (nm *NavMesh) addEdge(from, to NavNode) {
	if _, ok := nm.neighbors[from]; !ok {
		nm.neighbors[from] = []NavNode{}
	}
	nm.neighbors[from] = append(nm.neighbors[from], to)
}

func (nm *NavMesh) Neighbors(point NavNode) []NavNode {
	if neighbors, ok := nm.neighbors[point]; ok {
		return copyPointList(neighbors)
	}
	return []NavNode{}
}

func (nm *NavMesh) Cost(from, to geometry.Point) float64 {
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

	return length
}

func copyPointList(points []NavNode) []NavNode {
	return append([]NavNode{}, points...)
}
