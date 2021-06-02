package watcher

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	tektonv1alpha1 "github.com/vincent-pli/resource-watcher/pkg/apis/tekton/v1alpha1"
	handler "github.com/vincent-pli/resource-watcher/pkg/watcher/handlers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"

	// "knative.dev/eventing/pkg/kncloudevents"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/pkg/apis/duck"
	v1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Watcher struct {
	DiscoveryClient *discovery.DiscoveryClient
	DynamicClient   dynamic.Interface
	K8sClient       client.Client
	SwName          string
	SwNamespace     string
	Instance        *tektonv1alpha1.ResourceWatcher
	Log             logr.Logger
	Cache           cache.Cache
}

func (w Watcher) Start(stopCh <-chan struct{}) {
	// get ResourceWatcher cr
	instance := &tektonv1alpha1.ResourceWatcher{}
	nameNamespace := types.NamespacedName{
		Name:      w.SwName,
		Namespace: w.SwNamespace,
	}

	err := w.K8sClient.Get(context.TODO(), nameNamespace, instance)
	if err != nil {
		w.Log.Error(err, "get sink raise exception")
	}
	w.Instance = instance

	// create cache
	namespaces := []string{}
	for _, namespace := range instance.Spec.Namespaces {
		namespaces = append(namespaces, namespace)
	}
	// cache, err := cache.MultiNamespacedCacheBuilder(namespaces)(ctrl.GetConfigOrDie(), cache.Options{})
	// if err != nil {
	// 	w.Log.Error(err, "create cache raise exception")
	// 	return err
	// }
	// w.Cache = cache
	// w.Log.Info("Create cache watch on: ", "namespaces", namespaces)

	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(w.DynamicClient, 0)

	// prepare eventHandler, could be couldevent or k8sevent
	sink, err := w.getSinkURI(w.Instance.Spec.Sink, w.Instance.Namespace)
	if err != nil {
		w.Log.Error(err, "get sink raise exception")
	}

	// eventsClient, err := kncloudevents.NewDefaultClient(sink)
	eventsClient, err := cloudevents.NewClientHTTP()
	if err != nil {
		w.Log.Error(err, "failed to create client")
	}

	if err != nil {
		w.Log.Error(err, "creat cloudevent client raise exception")
	}

	for _, resource := range w.Instance.Spec.Resources {
		kind := resource.Kind
		var resourceStr string

		gv, err := schema.ParseGroupVersion(resource.APIVersion)
		if err != nil {
			w.Log.Error(err, "error parsing APIVersion")
			continue
		}
		// gvk := schema.GroupVersionKind{Kind: kind, Group: gv.Group, Version: gv.Version}

		preferredResources, err := w.DiscoveryClient.ServerResourcesForGroupVersion(gv.String())
		if err != nil {
			if discovery.IsGroupDiscoveryFailedError(err) {
				// w.Log.Warningf("failed to discover some groups: %v", err.(*discovery.ErrGroupDiscoveryFailed).Groups)
				fmt.Println("failed to discover some groups")
			} else {
				// w.Log.Warningf("failed to discover preferred resources: %v", err)
				fmt.Println("failed to discover preferred resources")
			}
		}

		for _, r := range preferredResources.APIResources {
			if r.Kind == kind && strings.Index(r.Name, "status") == -1 {
				resourceStr = r.Name
			}
		}

		gvr := schema.GroupVersionResource{Resource: resourceStr, Group: gv.Group, Version: gv.Version}
		// informer, err := w.Cache.GetInformerForKind(gvk)
		// if err != nil {
		// 	w.Log.Error(err, "cannot get informer by gvk")
		// 	continue
		// }

		couldEventHandler := &handler.CouldeventHandler{
			Resource: resource.NameSelector,
			Client:   eventsClient,
			Sink:     sink,
			Log:      w.Log,
		}
		// informer.AddEventHandler(couldEventHandler)
		informerFactory.ForResource(gvr).Informer().AddEventHandler(couldEventHandler)

		fmt.Println(gvr.String())

	}
	informerFactory.Start(stopCh)
	<-stopCh
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
