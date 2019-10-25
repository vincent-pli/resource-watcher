package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceWatcherSpec defines the desired state of ResourceWatcher
// +k8s:openapi-gen=true
type ResourceWatcherSpec struct {
	Namespaces []string `json:"namespaces"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Resources is the list of resources to watch
	Resources []ApiServerResource `json:"resources"`
	// ServiceAccountName is the name of the ServiceAccount to use to run this
	// source.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Sink is a reference to an object that will resolve to a domain name to use as the sink.
	// +optional
	Sink *corev1.ObjectReference `json:"sink,omitempty"`
}

// ResourceWatcherStatus defines the observed state of ResourceWatcher
// +k8s:openapi-gen=true
type ResourceWatcherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// SinkURI is the current active sink URI that has been configured for the ApiServerSource.
	// +optional
	SinkURI string `json:"sinkUri,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceWatcher is the Schema for the resourcewatchers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ResourceWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceWatcherSpec   `json:"spec,omitempty"`
	Status ResourceWatcherStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceWatcherList contains a list of ResourceWatcher
type ResourceWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceWatcher{}, &ResourceWatcherList{})
}

// ApiServerResource defines the resource to watch
type ApiServerResource struct {
	// API version of the resource to watch.
	APIVersion string `json:"apiVersion"`

	// Kind of the resource to watch.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`

	// LabelSelector restricts this source to objects with the selected labels
	// More info: http://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	LabelSelector metav1.LabelSelector `json:"labelSelector"`

	// ControllerSelector restricts this source to objects with a controlling owner reference of the specified kind.
	// Only apiVersion and kind are used. Both are optional.
	ControllerSelector metav1.OwnerReference `json:"controllerSelector"`

	// If true, send an event referencing the object controlling the resource
	Controller bool `json:"controller"`

	// NameSelector is the list of resource name watched
	NameSelector []string `json:"nameSelector"`
}
