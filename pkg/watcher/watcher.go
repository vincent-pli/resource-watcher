package watcher

import (
	"context"

	"github.com/go-logr/logr"
	tektonv1alpha1 "github.com/vincent-pli/resource-watcher/pkg/apis/tekton/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Watcher struct {
	K8sClient   client.Client
	SwName      string
	SwNamespace string
	Instance    *tektonv1alpha1.ResourceWatcher
	Log         logr.Logger
	Cache       cache.Cache
}

func (w *Watcher) New() error {
	instance := &tektonv1alpha1.ResourceWatcher{}
	nameNamespace := types.NamespacedName{
		Name:      w.SwName,
		Namespace: w.SwNamespace,
	}

	err := w.K8sClient.Get(context.TODO(), nameNamespace, instance)
	if err != nil {
		return err
	}
	w.Instance = instance

	return nil
}

func (w *Watcher) Start() error {

	for _, resource := range w.Instance.Spec.Resources {
		kind := resource.Kind
		gv, err := schema.ParseGroupVersion(resource.APIVersion)
		if err != nil {
			w.Log.Error(err, "error parsing APIVersion")
			continue
		}

		gvk := schema.GroupVersionKind{Kind: kind, Group: gv.Group, Version: gv.Version}
		informer, err := w.Cache.GetInformerForKind(gvk)
		if err != nil {
			w.Log.Error(err, "cannot get informer by gvk")
			continue
		}
		informer.AddEventHandler()

	}
	return nil
}
