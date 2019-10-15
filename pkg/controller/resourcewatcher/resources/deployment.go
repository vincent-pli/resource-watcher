/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	"fmt"

	sourcesv1alpha1 "github.com/vincent-pli/resource-watcher/pkg/apis/tekton/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeWatchDeploy(source *sourcesv1alpha1.ResourceWatcher, watcherImage string) *v1.Deployment {
	replicas := int32(1)
	labels := map[string]string{
		"eventing-source":      "resource-watcher",
		"eventing-source-name": source.Name,
	}
	env := []corev1.EnvVar{
		{
			Name:  "WATCHER_NAME",
			Value: source.Name,
		},
		{
			Name:  "WATCHER_NAMESPACE",
			Value: source.Namespace,
		},
	}

	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-deployment", source.Name),
			Namespace:    source.Namespace,
			Labels:       labels,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "true",
					},
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: source.Spec.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:  "receive-adapter",
							Image: watcherImage,
							Env:   env,
						},
					},
				},
			},
		},
	}
}
