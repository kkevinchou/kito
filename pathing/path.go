package pathing

// README
// Find path from a to b.
// a and b may or may not be an existing node
// need to recompute graph with a and b added
// compute a path
// delete a and b
// return path

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

func (planner *Planner) FindPath(startNode *Node) {

}
