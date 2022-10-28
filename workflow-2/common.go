package workflow_2

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

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
