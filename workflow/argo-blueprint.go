package workflow

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"fmt"
	"strings"

	"github.com/heimdalr/dag"
)

type ArgoBlueprint struct {
	Dag dag.DAG `json:"dag,omitempty"`
}

func NewArgoBlueprint() *ArgoBlueprint {
	return &ArgoBlueprint{
		Dag: *dag.NewDAG(),
	}
}

func (argoBlueprint *ArgoBlueprint) AddNode(id string, node *ContainerNode) {
	argoBlueprint.Dag.AddVertexByID(id, node)
}

func (workflow *ArgoBlueprint) AddEdge(srcId string, dstId string) {
	workflow.Dag.AddEdge(srcId, dstId)
}

func (workflow *ArgoBlueprint) GetNodes() []*ContainerNode {
	vertices := workflow.Dag.GetVertices()
	var nodes = []*ContainerNode{}

	for _, v := range vertices {
		nodes = append(nodes, v.(*ContainerNode))
	}

	return nodes
}

func (blueprint *ArgoBlueprint) Load(argoWorkflow argo_adapter.Workflow) {
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
		blueprint.AddNode(node_id, &ContainerNode{
			Node: Node{
				Id:     node_id,
				Status: argoWfNode.Phase,
				Inputs: inputs,
				Ouputs: outputs,
			},
			Container: Container{}, // need add info container
		})
	}
	// Add Edges
	for node_id, argoWfNode := range argoWorkflow.Status.Nodes {
		for _, child_id := range argoWfNode.Children {
			blueprint.AddEdge(node_id, child_id)
		}
	}

}

func (blueprint *ArgoBlueprint) Submit() *ArgoWorkflow {
	fmt.Println("Submit an argo blue print")
	return NewArgoWorkflow(blueprint)
}

func (workflow *ArgoBlueprint) String() string {
	var sb strings.Builder
	vertices := workflow.Dag.GetVertices()
	size := workflow.Dag.GetSize()

	visitor := &dagVisitor{Dag: &workflow.Dag}
	workflow.Dag.BFSWalk(visitor)

	// edges := workflow.Dag.()
	sb.WriteString(fmt.Sprintf("DAG Vertices: %d\n", len(vertices)))
	for _, vertice := range vertices {
		verticeNode := vertice.(*ContainerNode)
		sb.WriteString(fmt.Sprintf("- %s\n", verticeNode.Id))
	}
	sb.WriteString(fmt.Sprintf("DAG Edges: %d\n", size))
	sb.WriteString(visitor.EdgesDesciber)
	// return sb.String()
	return sb.String()
}
