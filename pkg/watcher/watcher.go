package watcher

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	tektonv1alpha1 "github.com/vincent-pli/resource-watcher/pkg/apis/tekton/v1alpha1"
	handler "github.com/vincent-pli/resource-watcher/pkg/watcher/handlers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/eventing/pkg/kncloudevents"
	"knative.dev/pkg/apis/duck"
	v1beta1 "knative.dev/pkg/apis/duck/v1beta1"
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

func (w Watcher) Start(stopCh <-chan struct{}) error {
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

	// prepare eventHandler, could be couldevent or k8sevent
	sink, err := w.getSinkURI(w.Instance.Spec.Sink, w.Instance.Namespace)
	if err != nil {
		w.Log.Error(err, "get sink raise exception")
		return err
	}

	eventsClient, err := kncloudevents.NewDefaultClient(sink)
	if err != nil {
		w.Log.Error(err, "creat cloudevent client raise exception")
		return err
	}

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

		couldEventHandler := &handler.CouldeventHandler{
			Resource: resource.NameSelector,
			Client:   eventsClient,
			Sink:     sink,
			Log:      w.Log,
		}
		informer.AddEventHandler(couldEventHandler)

	}
	w.Cache.Start(stopCh)

	return nil
}

// GetSinkURI retrieves the sink URI from the object referenced by the given
// ObjectReference.
func (w Watcher) getSinkURI(sink *corev1.ObjectReference, namespace string) (string, error) {
	if sink == nil {
		return "", fmt.Errorf("sink ref is nil")
	}

	if sink.Namespace == "" {
		sink.Namespace = namespace
	}

	objIdentifier := fmt.Sprintf("\"%s/%s\" (%s)", sink.Namespace, sink.Name, sink.GroupVersionKind())

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(sink.GroupVersionKind())
	err := w.K8sClient.Get(context.TODO(), client.ObjectKey{Name: sink.Name, Namespace: sink.Namespace}, u)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize sink %s: %v", objIdentifier, err)
	}

	t := v1beta1.AddressableType{}
	err = duck.FromUnstructured(u, &t)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize sink %s: %v", objIdentifier, err)
	}

	if t.Status.Address == nil {
		return "", fmt.Errorf("sink %s does not contain address", objIdentifier)
	}

	if t.Status.Address.URL == nil {
		return "", fmt.Errorf("sink %s contains an empty URL", objIdentifier)
	}

	return t.Status.Address.URL.String(), nil
}
