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
	ExampleGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Example",
	}
	ExampleResource = metav1.APIResource{
		Name:         "examples",
		SingularName: "example",
		Namespaced:   false,
		Kind:         ExampleGroupVersionKind.Kind,
	}

	ExampleGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "examples",
	}
)

func init() {
	resource.Put(ExampleGroupVersionResource)
}

func NewExample(namespace, name string, obj Example) *Example {
	obj.APIVersion, obj.Kind = ExampleGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Example `json:"items"`
}

type ExampleHandlerFunc func(key string, obj *Example) (runtime.Object, error)

type ExampleChangeHandlerFunc func(obj *Example) (runtime.Object, error)

type ExampleLister interface {
	List(namespace string, selector labels.Selector) (ret []*Example, err error)
	Get(namespace, name string) (*Example, error)
}

type ExampleController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ExampleLister
	AddHandler(ctx context.Context, name string, handler ExampleHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ExampleHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ExampleHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler ExampleHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ExampleInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Example) (*Example, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Example, error)
	Get(name string, opts metav1.GetOptions) (*Example, error)
	Update(*Example) (*Example, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ExampleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ExampleController
	AddHandler(ctx context.Context, name string, sync ExampleHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ExampleHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ExampleLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ExampleLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ExampleHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ExampleHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ExampleLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ExampleLifecycle)
}

type exampleLister struct {
	controller *exampleController
}

func (l *exampleLister) List(namespace string, selector labels.Selector) (ret []*Example, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Example))
	})
	return
}

func (l *exampleLister) Get(namespace, name string) (*Example, error) {
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
			Group:    ExampleGroupVersionKind.Group,
			Resource: "example",
		}, key)
	}
	return obj.(*Example), nil
}

type exampleController struct {
	controller.GenericController
}

func (c *exampleController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *exampleController) Lister() ExampleLister {
	return &exampleLister{
		controller: c,
	}
}

func (c *exampleController) AddHandler(ctx context.Context, name string, handler ExampleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Example); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *exampleController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler ExampleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Example); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *exampleController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ExampleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Example); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *exampleController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler ExampleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Example); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type exampleFactory struct {
}

func (c exampleFactory) Object() runtime.Object {
	return &Example{}
}

func (c exampleFactory) List() runtime.Object {
	return &ExampleList{}
}

func (s *exampleClient) Controller() ExampleController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.exampleControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ExampleGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &exampleController{
		GenericController: genericController,
	}

	s.client.exampleControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type exampleClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ExampleController
}

func (s *exampleClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *exampleClient) Create(o *Example) (*Example, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Example), err
}

func (s *exampleClient) Get(name string, opts metav1.GetOptions) (*Example, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Example), err
}

func (s *exampleClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Example, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Example), err
}

func (s *exampleClient) Update(o *Example) (*Example, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Example), err
}

func (s *exampleClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *exampleClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *exampleClient) List(opts metav1.ListOptions) (*ExampleList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ExampleList), err
}

func (s *exampleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *exampleClient) Patch(o *Example, patchType types.PatchType, data []byte, subresources ...string) (*Example, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*Example), err
}

func (s *exampleClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *exampleClient) AddHandler(ctx context.Context, name string, sync ExampleHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *exampleClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ExampleHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *exampleClient) AddLifecycle(ctx context.Context, name string, lifecycle ExampleLifecycle) {
	sync := NewExampleLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *exampleClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ExampleLifecycle) {
	sync := NewExampleLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *exampleClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ExampleHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *exampleClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ExampleHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *exampleClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ExampleLifecycle) {
	sync := NewExampleLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *exampleClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ExampleLifecycle) {
	sync := NewExampleLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type ExampleIndexer func(obj *Example) ([]string, error)

type ExampleClientCache interface {
	Get(namespace, name string) (*Example, error)
	List(namespace string, selector labels.Selector) ([]*Example, error)

	Index(name string, indexer ExampleIndexer)
	GetIndexed(name, key string) ([]*Example, error)
}

type ExampleClient interface {
	Create(*Example) (*Example, error)
	Get(namespace, name string, opts metav1.GetOptions) (*Example, error)
	Update(*Example) (*Example, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ExampleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ExampleClientCache

	OnCreate(ctx context.Context, name string, sync ExampleChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ExampleChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ExampleChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ExampleInterface
}

type exampleClientCache struct {
	client *exampleClient2
}

type exampleClient2 struct {
	iface      ExampleInterface
	controller ExampleController
}

func (n *exampleClient2) Interface() ExampleInterface {
	return n.iface
}

func (n *exampleClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *exampleClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *exampleClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *exampleClient2) Create(obj *Example) (*Example, error) {
	return n.iface.Create(obj)
}

func (n *exampleClient2) Get(namespace, name string, opts metav1.GetOptions) (*Example, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *exampleClient2) Update(obj *Example) (*Example, error) {
	return n.iface.Update(obj)
}

func (n *exampleClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *exampleClient2) List(namespace string, opts metav1.ListOptions) (*ExampleList, error) {
	return n.iface.List(opts)
}

func (n *exampleClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *exampleClientCache) Get(namespace, name string) (*Example, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *exampleClientCache) List(namespace string, selector labels.Selector) ([]*Example, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *exampleClient2) Cache() ExampleClientCache {
	n.loadController()
	return &exampleClientCache{
		client: n,
	}
}

func (n *exampleClient2) OnCreate(ctx context.Context, name string, sync ExampleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &exampleLifecycleDelegate{create: sync})
}

func (n *exampleClient2) OnChange(ctx context.Context, name string, sync ExampleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &exampleLifecycleDelegate{update: sync})
}

func (n *exampleClient2) OnRemove(ctx context.Context, name string, sync ExampleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &exampleLifecycleDelegate{remove: sync})
}

func (n *exampleClientCache) Index(name string, indexer ExampleIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*Example); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *exampleClientCache) GetIndexed(name, key string) ([]*Example, error) {
	var result []*Example
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*Example); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *exampleClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type exampleLifecycleDelegate struct {
	create ExampleChangeHandlerFunc
	update ExampleChangeHandlerFunc
	remove ExampleChangeHandlerFunc
}

func (n *exampleLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *exampleLifecycleDelegate) Create(obj *Example) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *exampleLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *exampleLifecycleDelegate) Remove(obj *Example) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *exampleLifecycleDelegate) Updated(obj *Example) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
