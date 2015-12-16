package pathing

import "github.com/kkevinchou/ant/util"

// Assumtions for pathfinding
// 1. Polygons are non overlapping - though they can share vertices
// 2. Polygons are convex

var nodeIdGenerator int = 0

type Node struct {
	Id int
	X  int
	Y  int
}

type Edge struct {
	srcNode *Node
	dstNode *Node
}

type Planner struct {
	nodes map[int]*Node
	edges map[*Node][]*Edge
}

func CreatePlanner(nodes []*Node, edges []*Edge) *Planner {
	planner := &Planner{}
	planner.nodes = make(map[int]*Node)
	planner.edges = make(map[*Node][]*Edge)

	for _, node := range nodes {
		planner.nodes[node.Id] = node
	}

	for _, edge := range edges {
		if _, ok := planner.edges[edge.srcNode]; !ok {
			planner.edges[edge.srcNode] = make([]*Edge, 0)
		}
		planner.edges[edge.srcNode] = append(planner.edges[edge.srcNode], edge)
	}

	return planner
}

func CreateEdge(srcNode, dstNode *Node) *Edge {
	return &Edge{
		srcNode: srcNode,
		dstNode: dstNode,
	}
}

func CreateNode(x int, y int) *Node {
	node := Node{
		Id: nodeIdGenerator,
		X:  x,
		Y:  y,
	}
	nodeIdGenerator++
	return &node
}

func (planner *Planner) FindPath(startNode *Node, goalNode *Node) {
	cameFrom := map[*Node]*Node{}
	cameFrom[startNode] = nil

	costSoFar := map[*Node]int{}
	costSoFar[startNode] = 0

	open := priorityqueue.New()
	open.Push(&Item{})
}
