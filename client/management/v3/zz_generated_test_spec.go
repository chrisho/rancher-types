package client

const (
	TestSpecType             = "testSpec"
	TestSpecFieldDisplayName = "displayName"
)

type TestSpec struct {
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
}
