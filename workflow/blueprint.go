package workflow

type Blueprint interface {
	AddNode()
	AddEdge()
	GetNodes()
	Submit() IWorkflow
	Load()
}

// func NewBlueprint() *Blueprint {
// 	return &Blueprint{
// 		Dag: *dag.NewDAG(),
// 	}
// }

// func (workflow *Blueprint) String() string {
// 	var sb strings.Builder
// 	vertices := workflow.Dag.GetVertices()
// 	size := workflow.Dag.GetSize()

// 	visitor := &dagVisitor{Dag: &workflow.Dag}
// 	workflow.Dag.BFSWalk(visitor)

// 	// edges := workflow.Dag.()
// 	sb.WriteString(fmt.Sprintf("DAG Vertices: %d\n", len(vertices)))
// 	for _, vertice := range vertices {
// 		verticeNode := vertice.(*Node)
// 		sb.WriteString(fmt.Sprintf("- %s\n", verticeNode.Id))
// 	}
// 	sb.WriteString(fmt.Sprintf("DAG Edges: %d\n", size))
// 	sb.WriteString(visitor.EdgesDesciber)
// 	// return sb.String()
// 	return sb.String()
// }

// func (workflow *Blueprint) AddNode(id string, node *Node) {
// 	// workflow.Nodes[id] = node
// 	workflow.Dag.AddVertexByID(id, node)
// }

// func (workflow *Blueprint) GetNodes() []*Node {
// 	vertices := workflow.Dag.GetVertices()
// 	var nodes = []*Node{}

// 	for _, v := range vertices {
// 		nodes = append(nodes, v.(*Node))
// 	}

// 	return nodes
// }
