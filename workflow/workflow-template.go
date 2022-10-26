package workflow

type Action struct {
	Node
}

type WorkflowTemplate struct {
	Actions []Action
}
