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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kedav1alpha1 "github.com/kedacore/keda/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MakeScaledObject(source *sourcesv1alpha1.ResourceWatcher) *corev1.Service {
	labels := map[string]string{
		"eventing-source":      "resource-watcher",
		"eventing-source-name": source.Name,
	}

	return &kedav1alpha1.{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-service", source.Name),
			Namespace:    source.Namespace,
			Labels:       labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       6000,
					TargetPort: intstr.FromInt(6000),
				},
			},
		},
	}
}
