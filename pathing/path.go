package pathing

import (
	"github.com/kkevinchou/ant/geometry"
	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/util"
)

const (
	floatEpsilon = float64(0.001) * float64(0.001)
)

type Planner struct {
	navmesh *NavMesh
}

func (p *Planner) SetNavMesh(navmesh *NavMesh) {
	p.navmesh = navmesh
}

func (p *Planner) FindPath(start geometry.Point, goal geometry.Point) []Node {
	roughPath := p.findNodePath(start, goal)
	if roughPath == nil {
		return nil
	}

	startNode, goalNode := Node{X: start.X, Y: start.Y}, Node{X: goal.X, Y: goal.Y}
	portals := findPortals(startNode, goalNode, roughPath)

	smoothedPath := smoothPath(portals)

	return smoothedPath
}

func findPortals(start Node, goal Node, nodes []Node) []Portal {
	portals := []Portal{CreatePortal(start, start)}
	prevPolygon := nodes[0].Polygon

	for _, node := range nodes {
		if node.Polygon != prevPolygon {
			portals = append(portals, prevPolygon.neighbors[node.Polygon])
			prevPolygon = node.Polygon
		}
	}

	portals = append(portals, CreatePortal(goal, goal))

	return portals
}

func (p *Planner) findNodePath(start geometry.Point, goal geometry.Point) []Node {
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
	pathNode := goalNode

	for {
		path = append(path, pathNode)
		if pathNode == startNode {
			break
		}
		pathNode = cameFrom[pathNode]
	}

	reversePath := make([]Node, len(path))
	for i := 0; i < len(path); i++ {
		reversePath[len(path)-1-i] = path[i]
	}

	return reversePath
}

func orderPortalNodes(portals []Portal) []Node {
	portalNodes := []Node{portals[0].Node1, portals[0].Node2}
	prevLeft := portals[0].Node1
	prevRight := portals[0].Node2

	for i := 1; i < len(portals); i++ {
		nextLeft := portals[i].Node1
		nextRight := portals[i].Node2

		leftVec := nextLeft.Vector().Sub(prevLeft.Vector())
		rightVec := nextRight.Vector().Sub(prevRight.Vector())

		if rightVec.Cross(leftVec) > 0 {
			nextLeft, nextRight = nextRight, nextLeft
		}
		// TODO: handle where they're == 0

		portalNodes = append(portalNodes, nextLeft)
		portalNodes = append(portalNodes, nextRight)

		prevLeft, prevRight = nextLeft, nextRight
	}

	return portalNodes
}

// Returns true if v is to left of reference
func vecOnLeft(reference, v vector.Vector) bool {
	return reference.Cross(v) < floatEpsilon
	// return reference.Cross(v) < 0
}

// Returns true if v is to the right of reference
func vecOnRight(reference, v vector.Vector) bool {
	return reference.Cross(v) > -1*floatEpsilon
	// return reference.Cross(v) > 0
}

func smoothPath(unorderedPortals []Portal) []Node {
	portalNodes := orderPortalNodes(unorderedPortals)

	// This algorithm was retrieved online but a confusing note:
	// lastValidRightIndex actually represent "left" index
	//
	// lastValidRightIndex represents the left index of the last valid
	// right index.  These indexes are used purely to reset the apex
	// at the correct point

	lastValidLeftIndex := 0
	lastValidRightIndex := 0

	apex := portalNodes[0]
	portalLeft := apex
	portalRight := apex

	contactNodes := []Node{apex}

	for i := 2; i < len(portalNodes); i += 2 {
		leftNode := portalNodes[i]
		rightNode := portalNodes[i+1]

		leftVec := leftNode.Vector().Sub(apex.Vector())
		rightVec := rightNode.Vector().Sub(apex.Vector())
		lastValidLeftVec := portalLeft.Vector().Sub(apex.Vector())
		lastValidRightVec := portalRight.Vector().Sub(apex.Vector())

		// Left side of funnel
		// The leftVec is to the right of lastValidLeftVec, so we
		// shrink the funnel
		if vecOnLeft(leftVec, lastValidLeftVec) {
			if (portalLeft == apex) || !vecOnRight(lastValidRightVec, leftVec) {
				portalLeft = leftNode
				lastValidLeftIndex = i
			} else {
				// If the new leftVec is to the right of the last valid
				// right vec, we set the new apex
				apex = portalRight
				portalLeft = apex
				if contactNodes[len(contactNodes)-1] != apex {
					contactNodes = append(contactNodes, apex)
				}

				lastValidLeftIndex = lastValidRightIndex
				i = lastValidRightIndex
				continue
			}
		}

		// Right side of funnel
		if vecOnRight(rightVec, lastValidRightVec) {
			if (portalRight == apex) || !vecOnLeft(lastValidLeftVec, rightVec) {
				portalRight = rightNode
				lastValidRightIndex = i
			} else {
				apex = portalLeft
				portalRight = apex
				if contactNodes[len(contactNodes)-1] != apex {
					contactNodes = append(contactNodes, apex)
				}

				lastValidRightIndex = lastValidLeftIndex
				i = lastValidLeftIndex
				continue
			}
		}
	}

	if contactNodes[len(contactNodes)-1] != portalNodes[len(portalNodes)-1] {
		contactNodes = append(contactNodes, portalNodes[len(portalNodes)-1])
	}

	return contactNodes
}
