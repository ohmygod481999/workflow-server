package workflow_2

type Blueprint interface {
	AddNode()
	AddEdge()
	GetNodes()
	Submit() Workflow
	Load()
	Save()
}
