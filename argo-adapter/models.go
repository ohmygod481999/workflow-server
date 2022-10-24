package argo_adapter

type MetaData map[string]interface{}

type Parameter struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	ValueFrom struct {
		Path string `json:"path"`
	} `json:"valueFrom"`
}

type Artifact struct {
	Name string `json:"name"`
	Path string `json:"path"`
	S3   struct {
		Key string
	} `json:"s3"`
}

type Arguments struct {
	Parameters []Parameter
}

type WorkflowTemplateRef struct {
	Name string `json:"name"`
}

type Template struct {
	Name      string    `json:"name,omitempty"`
	Inputs    Inputs    `json:"inputs,omitempty"`
	Outputs   Outputs   `json:"outputs,omitempty"`
	MetaData  MetaData  `json:"metaData,omitempty"`
	Container Container `json:"container,omitempty"`
}

type Spec struct {
	Entrypoint          string              `json:"entrypoint"`
	Arguments           Arguments           `json:"arguments"`
	WorkflowTemplateRef WorkflowTemplateRef `json:"workflowTemplateRef"`
}

type Inputs struct {
	Parameters []Parameter `json:"parameters"`
	Artifacts  []Artifact  `json:"artifacts"`
}

type Outputs struct {
	Parameters []Parameter `json:"parameters"`
	Artifacts  []Artifact  `json:"artifacts"`
	ExitCode   string      `json:"exitCode"`
}

type ResourcesDuration struct {
	Cpu    int `json:"cpu"`
	Memory int `json:"int"`
}

type Node struct {
	Id                string            `json:"id"`
	Name              string            `json:"name"`
	DisplayName       string            `json:"displayName"`
	Type              string            `json:"type"`
	TemplateName      string            `json:"TemplateName"`
	TemplateScope     string            `json:"TemplateScope"`
	Phase             string            `json:"phase"`
	BoundaryID        string            `json:"boundaryID"`
	StartedAt         string            `json:"startedAt"`
	FinishedAt        string            `json:"finishedAt"`
	EstimatedDuration int               `json:"estimatedDuration"`
	Progress          string            `json:"progress"`
	ResourcesDuration ResourcesDuration `json:"resourcesDuration"`
	Inputs            Inputs            `json:"inputs"`
	Outputs           Outputs           `json:"outputs"`
	Children          []string          `json:"children"`
	OutboundNodes     []string          `json:"outboundNodes"`
	HostNodeName      string            `json:"hostNodeName"`
}

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

type StoredTemplate struct {
	Name      string    `json:"name"`
	Inputs    Inputs    `json:"inputs"`
	Outputs   Outputs   `json:"outputs"`
	Metadata  MetaData  `json:"metadata"`
	Container Container `json:"container"`
}

type Condition struct {
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

type S3 struct {
	Enpoint         string `json:"enpoint,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
	Insecure        bool   `json:"insecure,omitempty"`
	AccessKeySecret struct {
		Name string `json:"name,omitempty"`
		Key  string `json:"key,omitempty"`
	} `json:"accessKeySecret,omitempty"`
	SecretKeySecret struct {
		Name string `json:"name,omitempty"`
		Key  string `json:"key,omitempty"`
	} `json:"secretKeySecret,omitempty"`
}

type ArtifactRepository struct {
	ArchiveLogs bool `json:"archiveLogs,omitempty"`
	S3          S3   `json:"s3,omitempty"`
}

type ArtifactRepositoryRef struct {
	ConfigMap          string             `json:"configMap,omitempty"`
	Key                string             `json:"key,omitempty"`
	Namespace          string             `json:"namespace,omitempty"`
	ArtifactRepository ArtifactRepository `json:"artifactRepository,omitempty"`
}

type Status struct {
	Phase                      string                    `json:"phase,omitempty"`
	StartedAt                  string                    `json:"startedAt,omitempty"`
	FinishedAt                 string                    `json:"finishedAt,omitempty"`
	EstimatedDuration          int                       `json:"estimatedDuration,omitempty"`
	Progress                   string                    `json:"progress,omitempty"`
	Nodes                      map[string]Node           `json:"nodes,omitempty"`
	StoredTemplates            map[string]StoredTemplate `json:"storedTemplates,omitempty"`
	Conditions                 []Condition               `json:"conditions,omitempty"`
	ResourcesDuration          ResourcesDuration         `json:"resourcesDuration,omitempty"`
	StoredWorkflowTemplateSpec Spec                      `json:"storedWorkflowTemplateSpec,omitempty"`
	ArtifactRepositoryRef      ArtifactRepositoryRef     `json:"artifactRepositoryRef,omitempty"`
}

type Workflow struct {
	MetaData MetaData `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

type Workflows struct {
	MetaData MetaData   `json:"metadata"`
	Items    []Workflow `json:"items"`
}

type WorkflowEvent struct {
	Result struct {
		Type   string   `json:"type,omitempty"`
		Object Workflow `json:"object,omitempty"`
	} `json:"result,omitempty"`
}
