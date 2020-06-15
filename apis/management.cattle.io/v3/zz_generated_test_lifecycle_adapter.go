package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

type TestLifecycle interface {
	Create(obj *Test) (runtime.Object, error)
	Remove(obj *Test) (runtime.Object, error)
	Updated(obj *Test) (runtime.Object, error)
}

type testLifecycleAdapter struct {
	lifecycle TestLifecycle
}

func (w *testLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *testLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *testLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Test))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *testLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Test))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *testLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Test))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTestLifecycleAdapter(name string, clusterScoped bool, client TestInterface, l TestLifecycle) TestHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TestGroupVersionResource)
	}
	adapter := &testLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Test) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
