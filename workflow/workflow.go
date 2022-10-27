package workflow

import "fmt"

type SubmitableWorkflow interface {
	Submit()
}

type Workflow struct {
	Blueprint *Blueprint
}

func (workflow *Workflow) Submit() {
	fmt.Println("Submit workflow")
}
