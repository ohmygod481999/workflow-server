package workflow

import (
	"fmt"

	"github.com/heimdalr/dag"
)

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
