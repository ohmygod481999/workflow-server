package workflow

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"fmt"
	"strings"

	"github.com/heimdalr/dag"
)

type Workflow struct {
	Dag dag.DAG `json:"dag,omitempty"`
}

func NewWorkflow() *Workflow {
	return &Workflow{
		Dag: *dag.NewDAG(),
	}
}

type dagVisitor struct {
	Dag           *dag.DAG
	EdgesDesciber string
}

func (pv *dagVisitor) Visit(v dag.Vertexer) {
	_, value := v.Vertex()
	node := value.(*Node)
	parents, _ := pv.Dag.GetParents(node.Id)
	for _, parent := range parents {
		pv.EdgesDesciber += fmt.Sprintf("%s -> %s\n", parent.(*Node).Id, node.Id)
	}
}

func (workflow *Workflow) String() string {
	var sb strings.Builder
	vertices := workflow.Dag.GetVertices()
	size := workflow.Dag.GetSize()

	visitor := &dagVisitor{Dag: &workflow.Dag}
	workflow.Dag.BFSWalk(visitor)

	// edges := workflow.Dag.()
	sb.WriteString(fmt.Sprintf("DAG Vertices: %d\n", len(vertices)))
	for _, vertice := range vertices {
		verticeNode := vertice.(*Node)
		sb.WriteString(fmt.Sprintf("- %s\n", verticeNode.Id))
	}
	sb.WriteString(fmt.Sprintf("DAG Edges: %d\n", size))
	sb.WriteString(visitor.EdgesDesciber)
	// return sb.String()
	return sb.String()
}

func (workflow *Workflow) AddNode(id string, node *Node) {
	// workflow.Nodes[id] = node
	workflow.Dag.AddVertexByID(id, node)
}

func (workflow *Workflow) AddEdge(srcId string, dstId string) {
	workflow.Dag.AddEdge(srcId, dstId)
}

func (workflow *Workflow) ReadFromArgoWorkflow(argoWorkflow argo_adapter.Workflow) {
	// Add Nodes
	for node_id, argoWfNode := range argoWorkflow.Status.Nodes {
		inputs := []Input{}
		for _, parameter := range argoWfNode.Inputs.Parameters {
			inputs = append(inputs, Input{
				Name:  parameter.Name,
				Value: []byte(parameter.Value),
			})
		}
		outputs := []Output{}
		for _, parameter := range argoWfNode.Outputs.Parameters {
			outputs = append(outputs, Output{
				Name:  parameter.Name,
				Value: []byte(parameter.Value),
			})
		}
		workflow.AddNode(node_id, &Node{
			Id:     node_id,
			Status: argoWfNode.Phase,
			Inputs: inputs,
			Ouputs: outputs,
		})
	}
	// Add Edges
	for node_id, argoWfNode := range argoWorkflow.Status.Nodes {
		for _, child_id := range argoWfNode.Children {
			workflow.AddEdge(node_id, child_id)
		}
	}

}

func (workflow *Workflow) GetNodes() []*Node {
	vertices := workflow.Dag.GetVertices()
	var nodes = []*Node{}

	for _, v := range vertices {
		nodes = append(nodes, v.(*Node))
	}

	return nodes
}
