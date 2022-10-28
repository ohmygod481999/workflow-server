package workflow_2

type Blueprint interface {
	AddNode()
	AddEdge()
	GetNodes() []*Template
	GetNode(id string) *Template
	Submit() Workflow
	Save() error
}
