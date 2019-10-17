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

package main

import (
	"fmt"
	"os"

	"github.com/vincent-pli/resource-watcher/pkg/apis"
	watcher "github.com/vincent-pli/resource-watcher/pkg/watcher"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	name      = "WATCHER_NAME"
	nameSpace = "WATCHER_NAMESPACE"
)

var (
	log    = logf.Log.WithName("cmd")
	scheme = runtime.NewScheme()
	// clients = client.Client
)

func main() {
	rwName, defined := os.LookupEnv(name)
	if !defined {
		err := fmt.Errorf("No environment variable found")
		log.Error(err, "required environment variable %q not defined", name)
		os.Exit(1)
	}

	rwNamespace, defined := os.LookupEnv(nameSpace)
	if !defined {
		err := fmt.Errorf("No environment variable found")
		log.Error(err, "required environment variable %q not defined", nameSpace)
		os.Exit(1)
	}
	// Setup Scheme for all resources
	if err := apis.AddToScheme(scheme); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	client, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme, Mapper: nil})
	if err != nil {
		log.Error(err, "exception raised when create client")
		os.Exit(1)
	}
	cache, err := cache.New(ctrl.GetConfigOrDie(), cache.Options{
		Scheme: scheme,
		Mapper: nil,
	})

	watcher := watcher.Watcher{
		K8sClient:   client,
		SwName:      rwName,
		SwNamespace: rwNamespace,
		Log:         log,
		Cache:       cache,
	}

	watcher.Start()

}
