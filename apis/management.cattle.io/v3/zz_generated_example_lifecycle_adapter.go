package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

type ExampleLifecycle interface {
	Create(obj *Example) (runtime.Object, error)
	Remove(obj *Example) (runtime.Object, error)
	Updated(obj *Example) (runtime.Object, error)
}

type exampleLifecycleAdapter struct {
	lifecycle ExampleLifecycle
}

func (w *exampleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *exampleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *exampleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Example))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *exampleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Example))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *exampleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Example))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewExampleLifecycleAdapter(name string, clusterScoped bool, client ExampleInterface, l ExampleLifecycle) ExampleHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ExampleGroupVersionResource)
	}
	adapter := &exampleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Example) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
