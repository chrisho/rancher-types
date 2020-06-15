package client

import (
	"github.com/rancher/norman/types"
)

const (
	TestType                 = "test"
	TestFieldAnnotations     = "annotations"
	TestFieldCreated         = "created"
	TestFieldCreatorID       = "creatorId"
	TestFieldDisplayName     = "displayName"
	TestFieldLabels          = "labels"
	TestFieldName            = "name"
	TestFieldOwnerReferences = "ownerReferences"
	TestFieldRemoved         = "removed"
	TestFieldUUID            = "uuid"
)

type Test struct {
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

type TestCollection struct {
	types.Collection
	Data   []Test `json:"data,omitempty"`
	client *TestClient
}

type TestClient struct {
	apiClient *Client
}

type TestOperations interface {
	List(opts *types.ListOpts) (*TestCollection, error)
	Create(opts *Test) (*Test, error)
	Update(existing *Test, updates interface{}) (*Test, error)
	Replace(existing *Test) (*Test, error)
	ByID(id string) (*Test, error)
	Delete(container *Test) error
}

func newTestClient(apiClient *Client) *TestClient {
	return &TestClient{
		apiClient: apiClient,
	}
}

func (c *TestClient) Create(container *Test) (*Test, error) {
	resp := &Test{}
	err := c.apiClient.Ops.DoCreate(TestType, container, resp)
	return resp, err
}

func (c *TestClient) Update(existing *Test, updates interface{}) (*Test, error) {
	resp := &Test{}
	err := c.apiClient.Ops.DoUpdate(TestType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *TestClient) Replace(obj *Test) (*Test, error) {
	resp := &Test{}
	err := c.apiClient.Ops.DoReplace(TestType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *TestClient) List(opts *types.ListOpts) (*TestCollection, error) {
	resp := &TestCollection{}
	err := c.apiClient.Ops.DoList(TestType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *TestCollection) Next() (*TestCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &TestCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *TestClient) ByID(id string) (*Test, error) {
	resp := &Test{}
	err := c.apiClient.Ops.DoByID(TestType, id, resp)
	return resp, err
}

func (c *TestClient) Delete(container *Test) error {
	return c.apiClient.Ops.DoResourceDelete(TestType, &container.Resource)
}
