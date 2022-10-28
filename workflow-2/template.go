package workflow_2

import (
	argo_adapter "callbot/workflow/argo-adapter"
	"fmt"
)

type TemplateArgoRef struct {
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
}

type Template struct {
	Id string `json:"id,omitempty"`

	Inputs  []Input  `json:"inputs,omitempty"`
	Configs []Config `json:"configs,omitempty"`
	Ouputs  []Output `json:"ouputs,omitempty"`

	TemplateArgoRef TemplateArgoRef `json:"templateArgoRef,omitempty"`
	Container       *Container      `json:"container,omitempty"`
	ParentIds       []string        `json:"parentIds,omitempty"`
}

type TemplateGroup struct {
	Id string `json:"id,omitempty"`

	Templates []Template `json:"templates,omitempty"`
}

func LoadTemplateGroupFromArgo(argoWorkflowTemplate *argo_adapter.WorkflowTemplate) *TemplateGroup {
	fmt.Println(argoWorkflowTemplate.MetaData)

	id := argoWorkflowTemplate.MetaData["name"]

	var templates []Template

	for _, argotemplate := range argoWorkflowTemplate.Spec.Templates {
		var container *Container
		// if argotemplate.Container.Image != "" {
		// 	container = &Container{
		// 		Name:    argotemplate.Container.Name,
		// 		Image:   argotemplate.Container.Image,
		// 		Command: argotemplate.Container.Command,
		// 		Args:    argotemplate.Container.Args,
		// 	}
		// } else {
		// 	container = nil
		// }
		var inputs []Input
		for _, param := range argotemplate.Inputs.Parameters {
			inputs = append(inputs, Input{
				Name:  param.Name,
				Value: []byte(param.Value),
			})
		}
		var outputs []Output
		for _, param := range argotemplate.Outputs.Parameters {
			outputs = append(outputs, Output{
				Name:  param.Name,
				Value: []byte(param.Value),
			})
		}
		templates = append(templates, Template{
			Id:        argotemplate.Name,
			Inputs:    inputs,
			Configs:   []Config{},
			Ouputs:    outputs,
			Container: container,
			TemplateArgoRef: TemplateArgoRef{
				Name:     id.(string),
				Template: argotemplate.Name,
			},
		})
	}

	return &TemplateGroup{
		Id:        id.(string),
		Templates: templates,
	}
}

func (template *Template) String() string {
	// result := fmt.Sprintf("{\n\tId: %s;\n\tStatus: %s\n\n\tInputs:}\n", node.Id, node.Status)

	result := "{"
	result += fmt.Sprintf("\n\tId: %s,", template.Id)
	result += fmt.Sprintf("\n\tRef: %s/%s,", template.TemplateArgoRef.Name, template.TemplateArgoRef.Template)
	result += fmt.Sprintf("\n\tInputs:")
	for _, input := range template.Inputs {
		result += fmt.Sprintf("\n\t\t- %s: %s", input.Name, input.Value)
	}
	result += fmt.Sprintf("\n\tParents:")
	for _, parentId := range template.ParentIds {
		result += fmt.Sprintf("\n\t\t- %s", parentId)
	}
	result += fmt.Sprintf("\n\tOuputs:")
	for _, output := range template.Ouputs {
		result += fmt.Sprintf("\n\t\t- %s: %s", output.Name, output.Value)
	}
	result += "\n}"
	return result

}
