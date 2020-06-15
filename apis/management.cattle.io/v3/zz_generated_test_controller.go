package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	TestGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Test",
	}
	TestResource = metav1.APIResource{
		Name:         "tests",
		SingularName: "test",
		Namespaced:   false,
		Kind:         TestGroupVersionKind.Kind,
	}

	TestGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "tests",
	}
)

func init() {
	resource.Put(TestGroupVersionResource)
}

func NewTest(namespace, name string, obj Test) *Test {
	obj.APIVersion, obj.Kind = TestGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Test `json:"items"`
}

type TestHandlerFunc func(key string, obj *Test) (runtime.Object, error)

type TestChangeHandlerFunc func(obj *Test) (runtime.Object, error)

type TestLister interface {
	List(namespace string, selector labels.Selector) (ret []*Test, err error)
	Get(namespace, name string) (*Test, error)
}

type TestController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() TestLister
	AddHandler(ctx context.Context, name string, handler TestHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync TestHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler TestHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler TestHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type TestInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Test) (*Test, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Test, error)
	Get(name string, opts metav1.GetOptions) (*Test, error)
	Update(*Test) (*Test, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*TestList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() TestController
	AddHandler(ctx context.Context, name string, sync TestHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync TestHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle TestLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle TestLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync TestHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync TestHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle TestLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle TestLifecycle)
}

type testLister struct {
	controller *testController
}

func (l *testLister) List(namespace string, selector labels.Selector) (ret []*Test, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Test))
	})
	return
}

func (l *testLister) Get(namespace, name string) (*Test, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    TestGroupVersionKind.Group,
			Resource: "test",
		}, key)
	}
	return obj.(*Test), nil
}

type testController struct {
	controller.GenericController
}

func (c *testController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *testController) Lister() TestLister {
	return &testLister{
		controller: c,
	}
}

func (c *testController) AddHandler(ctx context.Context, name string, handler TestHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Test); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *testController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler TestHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Test); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *testController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler TestHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Test); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *testController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler TestHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Test); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type testFactory struct {
}

func (c testFactory) Object() runtime.Object {
	return &Test{}
}

func (c testFactory) List() runtime.Object {
	return &TestList{}
}

func (s *testClient) Controller() TestController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.testControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(TestGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &testController{
		GenericController: genericController,
	}

	s.client.testControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type testClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   TestController
}

func (s *testClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *testClient) Create(o *Test) (*Test, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Test), err
}

func (s *testClient) Get(name string, opts metav1.GetOptions) (*Test, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Test), err
}

func (s *testClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Test, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Test), err
}

func (s *testClient) Update(o *Test) (*Test, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Test), err
}

func (s *testClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *testClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *testClient) List(opts metav1.ListOptions) (*TestList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*TestList), err
}

func (s *testClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *testClient) Patch(o *Test, patchType types.PatchType, data []byte, subresources ...string) (*Test, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*Test), err
}

func (s *testClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *testClient) AddHandler(ctx context.Context, name string, sync TestHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *testClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync TestHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *testClient) AddLifecycle(ctx context.Context, name string, lifecycle TestLifecycle) {
	sync := NewTestLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *testClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle TestLifecycle) {
	sync := NewTestLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *testClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync TestHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *testClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync TestHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *testClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle TestLifecycle) {
	sync := NewTestLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *testClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle TestLifecycle) {
	sync := NewTestLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type TestIndexer func(obj *Test) ([]string, error)

type TestClientCache interface {
	Get(namespace, name string) (*Test, error)
	List(namespace string, selector labels.Selector) ([]*Test, error)

	Index(name string, indexer TestIndexer)
	GetIndexed(name, key string) ([]*Test, error)
}

type TestClient interface {
	Create(*Test) (*Test, error)
	Get(namespace, name string, opts metav1.GetOptions) (*Test, error)
	Update(*Test) (*Test, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*TestList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() TestClientCache

	OnCreate(ctx context.Context, name string, sync TestChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync TestChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync TestChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() TestInterface
}

type testClientCache struct {
	client *testClient2
}

type testClient2 struct {
	iface      TestInterface
	controller TestController
}

func (n *testClient2) Interface() TestInterface {
	return n.iface
}

func (n *testClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *testClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *testClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *testClient2) Create(obj *Test) (*Test, error) {
	return n.iface.Create(obj)
}

func (n *testClient2) Get(namespace, name string, opts metav1.GetOptions) (*Test, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *testClient2) Update(obj *Test) (*Test, error) {
	return n.iface.Update(obj)
}

func (n *testClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *testClient2) List(namespace string, opts metav1.ListOptions) (*TestList, error) {
	return n.iface.List(opts)
}

func (n *testClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *testClientCache) Get(namespace, name string) (*Test, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *testClientCache) List(namespace string, selector labels.Selector) ([]*Test, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *testClient2) Cache() TestClientCache {
	n.loadController()
	return &testClientCache{
		client: n,
	}
}

func (n *testClient2) OnCreate(ctx context.Context, name string, sync TestChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &testLifecycleDelegate{create: sync})
}

func (n *testClient2) OnChange(ctx context.Context, name string, sync TestChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &testLifecycleDelegate{update: sync})
}

func (n *testClient2) OnRemove(ctx context.Context, name string, sync TestChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &testLifecycleDelegate{remove: sync})
}

func (n *testClientCache) Index(name string, indexer TestIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*Test); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *testClientCache) GetIndexed(name, key string) ([]*Test, error) {
	var result []*Test
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*Test); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *testClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type testLifecycleDelegate struct {
	create TestChangeHandlerFunc
	update TestChangeHandlerFunc
	remove TestChangeHandlerFunc
}

func (n *testLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *testLifecycleDelegate) Create(obj *Test) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *testLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *testLifecycleDelegate) Remove(obj *Test) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *testLifecycleDelegate) Updated(obj *Test) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
