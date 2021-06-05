package resourcewatcher

import (
	"context"
	"fmt"
	"os"

	tektonv1alpha1 "github.com/vincent-pli/resource-watcher/pkg/apis/tekton/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	resources "github.com/vincent-pli/resource-watcher/pkg/controller/resourcewatcher/resources"
)

var log = logf.Log.WithName("controller_resourcewatcher")

const (
	watcherImageEnvVar = "WATCH_IMAGE"
	finalizerName      = "resourcewatchers.tekton.dev/finalizer"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ResourceWatcher Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	watcherImage, defined := os.LookupEnv(watcherImageEnvVar)
	if !defined {
		err := fmt.Errorf("No environment variable found")
		log.Error(err, "required environment variable %q not defined", watcherImageEnvVar)
		return nil
	}
	return &ReconcileResourceWatcher{
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		watcherImage: watcherImage,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("resourcewatcher-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ResourceWatcher
	err = c.Watch(&source.Kind{Type: &tektonv1alpha1.ResourceWatcher{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ResourceWatcher
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tektonv1alpha1.ResourceWatcher{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileResourceWatcher implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileResourceWatcher{}

// ReconcileResourceWatcher reconciles a ResourceWatcher object
type ReconcileResourceWatcher struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	scheme       *runtime.Scheme
	watcherImage string
}

// Reconcile reads that state of the cluster for a ResourceWatcher object and makes changes based on the state read
// and what is in the ResourceWatcher.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileResourceWatcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ResourceWatcher")

	// Fetch the ResourceWatcher instance
	instance := &tektonv1alpha1.ResourceWatcher{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(instance.GetFinalizers(), finalizerName) {
			addFinalizer(instance, finalizerName)
			if err := r.client.Update(context.TODO(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(instance.GetFinalizers(), finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.removeClusterrolebinding(instance); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			removeFinalizer(instance, finalizerName)
			if err := r.client.Update(context.TODO(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return reconcile.Result{}, nil
	}

	// Define a new Rolebinding object
	rolebinding, err := resources.MakeRolebinding(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Set JobFlow instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, rolebinding, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Rolebinding already exists
	foundRolebinding := &rbacv1.ClusterRoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: rolebinding.Name, Namespace: rolebinding.Namespace}, foundRolebinding)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Rolebinding", "Rolebinding.Namespace", rolebinding.Namespace, "Rolebinding.Name", rolebinding.Name)
		err = r.client.Create(context.TODO(), rolebinding)
		if err != nil {
			reqLogger.Error(err, "Create Rolebinding raise exception.")
			return reconcile.Result{}, err
		}
	}

	// Define a new Pod object
	deploy := resources.MakeWatchDeploy(instance, r.watcherImage)

	// Set ResourceWatcher instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &v1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.client.Create(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileResourceWatcher) removeClusterrolebinding(o *tektonv1alpha1.ResourceWatcher) error {
	name := resources.GetRolebindingName(o)

	// Check if this Rolebinding already exists
	foundRolebinding := &rbacv1.ClusterRoleBinding{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ""}, foundRolebinding)
	if err == nil {
		err = r.client.Delete(context.TODO(), foundRolebinding)
	}

	if err != nil {
		return err
	}

	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// AddFinalizer accepts an Object and adds the provided finalizer if not present.
func addFinalizer(o *tektonv1alpha1.ResourceWatcher, finalizer string) {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return
		}
	}
	o.SetFinalizers(append(f, finalizer))
}

// RemoveFinalizer accepts an Object and removes the provided finalizer if present.
func removeFinalizer(o *tektonv1alpha1.ResourceWatcher, finalizer string) {
	f := o.GetFinalizers()
	for i := 0; i < len(f); i++ {
		if f[i] == finalizer {
			f = append(f[:i], f[i+1:]...)
			i--
		}
	}
	o.SetFinalizers(f)
}
