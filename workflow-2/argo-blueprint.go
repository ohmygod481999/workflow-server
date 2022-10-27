package workflow_2

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"fmt"
	"strings"

	"github.com/heimdalr/dag"
)

type ArgoBlueprint struct {
	Dag     dag.DAG `json:"dag,omitempty"`
	Inputs  []Input
	Outputs []Output
}

func NewArgoBlueprint() *ArgoBlueprint {
	return &ArgoBlueprint{
		Dag: *dag.NewDAG(),
	}
}

func (argoBlueprint *ArgoBlueprint) AddNode(id string, node *Template) {
	argoBlueprint.Dag.AddVertexByID(id, node)
}

func (workflow *ArgoBlueprint) AddEdge(srcId string, dstId string) {
	workflow.Dag.AddEdge(srcId, dstId)
}

func (workflow *ArgoBlueprint) GetNodes() []*Template {
	vertices := workflow.Dag.GetVertices()
	var nodes = []*Template{}

	for _, v := range vertices {
		nodes = append(nodes, v.(*Template))
	}

	return nodes
}

func (blueprint *ArgoBlueprint) Save() {
	fmt.Println("Saved blueprint")
}

func (blueprint *ArgoBlueprint) Submit() Workflow {
	fmt.Println("Saved blueprint")

	return NewArgoWorkflow(blueprint)
}

func (argoBlueprint *ArgoBlueprint) String() string {
	var sb strings.Builder
	vertices := argoBlueprint.Dag.GetVertices()
	size := argoBlueprint.Dag.GetSize()

	visitor := &dagVisitor{Dag: &argoBlueprint.Dag}
	argoBlueprint.Dag.BFSWalk(visitor)

	sb.WriteString(fmt.Sprintf("DAG Vertices: %d\n", len(vertices)))
	for _, vertice := range vertices {
		verticeNode := vertice.(*Template)
		sb.WriteString(fmt.Sprintf("- %s\n", verticeNode.Id))
	}
	sb.WriteString(fmt.Sprintf("DAG Edges: %d\n", size))
	sb.WriteString(visitor.EdgesDesciber)

	return sb.String()
}

func LoadBlueprint(argoWorkflow argo_adapter.Workflow) *ArgoBlueprint {
	// Add Nodes
	blueprint := NewArgoBlueprint()

	entrypoint := argoWorkflow.Spec.Entrypoint

	var entryTemplate *argo_adapter.Template = nil

	for _, argoTemplate := range argoWorkflow.Spec.Templates {
		if argoTemplate.Name == entrypoint {
			entryTemplate = &argoTemplate
		}
	}

	for _, task := range entryTemplate.Dag.Tasks {
		inputs := []Input{}
		for _, parameter := range task.Arguments.Parameters {
			inputs = append(inputs, Input{
				Name:  parameter.Name,
				Value: []byte(parameter.Value),
			})
		}
		outputs := []Output{}
		blueprint.AddNode(task.Name, &Template{
			Id:        task.Name,
			Inputs:    inputs,
			Ouputs:    outputs,
			Configs:   []Config{},
			Container: &Container{},
			TemplateArgoRef: TemplateArgoRef{
				Name:     task.TemplateRef.Name,
				Template: task.TemplateRef.Template,
			},
		})
		if task.Depends != "" {
			parents := strings.Split(task.Depends, "&&")
			for _, parent := range parents {
				blueprint.AddEdge(parent, task.Name)
			}
		}
	}

	return blueprint
}

// func Load(argoWorkflow argo_adapter.Workflow) *ArgoBlueprint {
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

// 	return blueprint
// }
