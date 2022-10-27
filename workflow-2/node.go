package workflow_2

import "fmt"

type Input struct {
	Name  string
	Value []byte
}

type Config struct {
	Name  string
	Value []byte
}

type Output struct {
	Name  string
	Value []byte
}

type Node struct {
	Id     string
	Status string

	Inputs []Input
	Ouputs []Output
}

func (node *Node) String() string {
	// result := fmt.Sprintf("{\n\tId: %s;\n\tStatus: %s\n\n\tInputs:}\n", node.Id, node.Status)

	result := "{"
	result += fmt.Sprintf("\n\tId: %s,", node.Id)
	result += fmt.Sprintf("\n\tStatus: %s,", node.Status)
	result += fmt.Sprintf("\n\tInputs:")
	for _, input := range node.Inputs {
		result += fmt.Sprintf("\n\t\t- %s: %s", input.Name, input.Value)
	}
	result += fmt.Sprintf("\n\tOuputs:")
	for _, output := range node.Ouputs {
		result += fmt.Sprintf("\n\t\t- %s: %s", output.Name, output.Value)
	}
	result += "\n}"
	return result

}
