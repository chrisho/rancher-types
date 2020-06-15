package client

const (
	ExampleSpecType             = "exampleSpec"
	ExampleSpecFieldDisplayName = "displayName"
)

type ExampleSpec struct {
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
}
