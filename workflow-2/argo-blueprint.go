package workflow_2

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	dstNode := workflow.GetNode(dstId)
	dstNode.ParentIds = append(dstNode.ParentIds, srcId)
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

func (workflow *ArgoBlueprint) GetNode(id string) *Template {
	nodes := workflow.GetNodes()
	var result *Template
	for _, node := range nodes {
		if node.Id == id {
			result = node
		}
	}
	return result
}

func (blueprint *ArgoBlueprint) Save() error {
	var workflowParams []argo_adapter.Parameter
	for _, param := range blueprint.Inputs {
		workflowParams = append(workflowParams, argo_adapter.Parameter{
			Name:  param.Name,
			Value: string(param.Value),
		})
	}

	var tasks []argo_adapter.Task

	// Loop all node in blueprint
	for _, node := range blueprint.GetNodes() {
		// Node input
		var inputParams []argo_adapter.Parameter
		for _, input := range node.Inputs {
			inputParams = append(inputParams, argo_adapter.Parameter{
				Name:  input.Name,
				Value: string(input.Value),
			})
		}
		// Node output
		var outputParams []argo_adapter.Parameter
		for _, output := range node.Ouputs {
			outputParams = append(outputParams, argo_adapter.Parameter{
				Name:  output.Name,
				Value: string(output.Value),
			})
		}

		// append task to tasks
		tasks = append(tasks, argo_adapter.Task{
			Name: node.Id,
			Arguments: argo_adapter.Arguments{
				Parameters: inputParams,
			},
			TemplateRef: argo_adapter.WorkflowTemplateRef{
				Name:     node.TemplateArgoRef.Name,
				Template: node.TemplateArgoRef.Template,
			},
			Depends: strings.Join(node.ParentIds, " && "),
		})
	}

	template := argo_adapter.DagTemplate{
		Template: argo_adapter.Template{
			Name:     "main",
			Inputs:   argo_adapter.Inputs{},
			MetaData: argo_adapter.MetaData{},
		},
		Dag: argo_adapter.Dag{
			Tasks: tasks,
		},
	}

	persisWorkflow := argo_adapter.Workflow{
		MetaData: argo_adapter.MetaData{
			"name":      fmt.Sprintf("test-%d", rand.Intn(1000)),
			"namepsace": "argo",
		},
		Spec: argo_adapter.Spec{
			Arguments: argo_adapter.Arguments{
				Parameters: workflowParams,
			},
			Entrypoint: "main",
			Templates:  []argo_adapter.DagTemplate{template},
		},
	}

	str, _ := json.MarshalIndent(persisWorkflow, "", "\t")

	filename := "saved-blueprint.json"
	err := ioutil.WriteFile(filename, str, 0644)

	fmt.Println("Saved blueprint to " + filename)

	return err
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

	var entryTemplate *argo_adapter.DagTemplate = nil

	for _, argoTemplate := range argoWorkflow.Spec.Templates {
		if argoTemplate.Name == entrypoint {
			entryTemplate = &argoTemplate
		}
	}

	var inputs []Input
	for _, input := range argoWorkflow.Spec.Arguments.Parameters {
		inputs = append(inputs, Input{
			Name:  input.Name,
			Value: []byte(input.Value),
		})
	}

	blueprint.Inputs = inputs

	for _, task := range entryTemplate.Dag.Tasks {
		inputs := []Input{}
		for _, parameter := range task.Arguments.Parameters {
			inputs = append(inputs, Input{
				Name:  parameter.Name,
				Value: []byte(parameter.Value),
			})
		}
		outputs := []Output{}
		template := Template{
			Id:        task.Name,
			Inputs:    inputs,
			Ouputs:    outputs,
			Configs:   []Config{},
			Container: &Container{},
			TemplateArgoRef: TemplateArgoRef{
				Name:     task.TemplateRef.Name,
				Template: task.TemplateRef.Template,
			},
			ParentIds: []string{},
		}
		blueprint.AddNode(task.Name, &template)

		// If task has "Depends" => Add edge from this task to those depends node
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
