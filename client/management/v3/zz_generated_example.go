package client

import (
	"github.com/rancher/norman/types"
)

const (
	ExampleType                 = "example"
	ExampleFieldAnnotations     = "annotations"
	ExampleFieldCreated         = "created"
	ExampleFieldCreatorID       = "creatorId"
	ExampleFieldDisplayName     = "displayName"
	ExampleFieldLabels          = "labels"
	ExampleFieldName            = "name"
	ExampleFieldOwnerReferences = "ownerReferences"
	ExampleFieldRemoved         = "removed"
	ExampleFieldUUID            = "uuid"
)

type Example struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	DisplayName     string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ExampleCollection struct {
	types.Collection
	Data   []Example `json:"data,omitempty"`
	client *ExampleClient
}

type ExampleClient struct {
	apiClient *Client
}

type ExampleOperations interface {
	List(opts *types.ListOpts) (*ExampleCollection, error)
	Create(opts *Example) (*Example, error)
	Update(existing *Example, updates interface{}) (*Example, error)
	Replace(existing *Example) (*Example, error)
	ByID(id string) (*Example, error)
	Delete(container *Example) error
}

func newExampleClient(apiClient *Client) *ExampleClient {
	return &ExampleClient{
		apiClient: apiClient,
	}
}

func (c *ExampleClient) Create(container *Example) (*Example, error) {
	resp := &Example{}
	err := c.apiClient.Ops.DoCreate(ExampleType, container, resp)
	return resp, err
}

func (c *ExampleClient) Update(existing *Example, updates interface{}) (*Example, error) {
	resp := &Example{}
	err := c.apiClient.Ops.DoUpdate(ExampleType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ExampleClient) Replace(obj *Example) (*Example, error) {
	resp := &Example{}
	err := c.apiClient.Ops.DoReplace(ExampleType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ExampleClient) List(opts *types.ListOpts) (*ExampleCollection, error) {
	resp := &ExampleCollection{}
	err := c.apiClient.Ops.DoList(ExampleType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ExampleCollection) Next() (*ExampleCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ExampleCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ExampleClient) ByID(id string) (*Example, error) {
	resp := &Example{}
	err := c.apiClient.Ops.DoByID(ExampleType, id, resp)
	return resp, err
}

func (c *ExampleClient) Delete(container *Example) error {
	return c.apiClient.Ops.DoResourceDelete(ExampleType, &container.Resource)
}
