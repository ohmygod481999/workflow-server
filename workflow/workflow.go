package workflow

import "fmt"

type IWorkflow interface {
	Watch()
}

type ArgoWorkflow struct {
	Blueprint *ArgoBlueprint
}

func NewArgoWorkflow(blueprint *ArgoBlueprint) *ArgoWorkflow {
	return &ArgoWorkflow{
		Blueprint: blueprint,
	}
}

func (workflow *ArgoWorkflow) Watch() {
	fmt.Println("Watch workflow")
}
