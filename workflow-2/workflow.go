package workflow_2

import (
	"fmt"
)

type Workflow interface {
	Watch()
}

type ArgoWorkflow struct {
	ArgoBlueprint ArgoBlueprint
}

func NewArgoWorkflow(argoBlueprint *ArgoBlueprint) *ArgoWorkflow {
	return &ArgoWorkflow{
		ArgoBlueprint: *argoBlueprint,
	}
}

// func LoadWorkflow(argoWorkflow argo_adapter.Workflow) *ArgoWorkflow {
// 	// Add Nodes
// 	blueprint := NewArgoBlueprint()

// 	for node_id, argoWfNode := range argoWorkflow.Status.Nodes {
// 		inputs := []Input{}
// 		for _, parameter := range argoWfNode.Inputs.Parameters {
// 			inputs = append(inputs, Input{
// 				Name:  parameter.Name,
// 				Value: []byte(parameter.Value),
// 			})
// 		}
// 		outputs := []Output{}
// 		for _, parameter := range argoWfNode.Outputs.Parameters {
// 			outputs = append(outputs, Output{
// 				Name:  parameter.Name,
// 				Value: []byte(parameter.Value),
// 			})
// 		}
// 		blueprint.AddNode(node_id, &ContainerNode{
// 			Node: Node{
// 				Id:     node_id,
// 				Status: argoWfNode.Phase,
// 				Inputs: inputs,
// 				Ouputs: outputs,
// 			},
// 			Container: Container{}, // need add info container
// 		})
// 	}
// 	// Add Edges
// 	for node_id, argoWfNode := range argoWorkflow.Status.Nodes {
// 		for _, child_id := range argoWfNode.Children {
// 			blueprint.AddEdge(node_id, child_id)
// 		}
// 	}

// 	return NewArgoWorkflow(blueprint)
// }

func (argoWorkflow *ArgoWorkflow) Watch() {
	fmt.Println("Watch")
}
