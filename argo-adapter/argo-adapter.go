package argo_adapter

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	sse "github.com/r3labs/sse/v2"
)

type Config struct {
	ArgoUri string
}

type WorkflowTemplateSubmitBody struct {
	Namespace     string `json:"namespace"`
	ResourceKind  string `json:"resourceKind"`
	ResourceName  string `json:"resourceName"`
	SubmitOptions struct {
		Entrypoint string   `json:"entryPoint"`
		Parameters []string `json:"parameters"`
	} `json:"submitOptions"`
}

type ArgoAdapter struct {
	client resty.Client
	uri    string
}

func New(cfg Config) *ArgoAdapter {
	client := resty.New()
	client.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	return &ArgoAdapter{
		client: *client,
		uri:    cfg.ArgoUri,
	}
}

func (argoAdapter *ArgoAdapter) GetWorkflows() *Workflows {
	url := fmt.Sprintf("%s/api/v1/workflows/argo", argoAdapter.uri)

	res, err := argoAdapter.client.R().Get(url)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil
	}

	var wf Workflows
	err = json.Unmarshal(res.Body(), &wf)

	if err != nil {
		fmt.Println(err)
	}

	return &wf
}

func (argoAdapter *ArgoAdapter) GetWorkflowTemplates() *WorkflowTemplates {
	url := fmt.Sprintf("%s/api/v1/workflow-templates/argo", argoAdapter.uri)

	res, err := argoAdapter.client.R().Get(url)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil
	}

	var wfTemplates WorkflowTemplates
	err = json.Unmarshal(res.Body(), &wfTemplates)

	if err != nil {
		fmt.Println(err)
	}

	return &wfTemplates

}

func (argoAdapter *ArgoAdapter) SubmitWorkflowTemplate(body WorkflowTemplateSubmitBody) (*Workflow, error) {
	url := fmt.Sprintf("%s/api/v1/workflows/argo/submit", argoAdapter.uri)
	res, err := argoAdapter.client.R().SetBody(body).Post(url)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	var wf Workflow
	err = json.Unmarshal(res.Body(), &wf)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	return &wf, nil
}

func (argoAdapter *ArgoAdapter) ListenWorkflowEvent(name string, namespace string, message chan WorkflowEvent) {
	url := fmt.Sprintf("%s/api/v1/workflow-events/argo?listOptions.fieldSelector=metadata.namespace=%s,metadata.name=%s", argoAdapter.uri, namespace, name)
	client := sse.NewClient(url)

	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		var wfEvent WorkflowEvent
		json.Unmarshal(msg.Data, &wfEvent)
		message <- wfEvent
	})
}
