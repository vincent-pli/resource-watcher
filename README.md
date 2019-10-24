# Resource-watcher

The project is used for watch on resources of k8s,
when anything changed on the watched resource (add, update, delete) will take action:
- Send `Cloudevent` to specific `url`
- Create k8s event
- TODO

The `action` should be pluginable, user could implements what they want, for example: send mail...

It's k8s native and implements by a k8s controller.


## Development Prerequisites
1. [`go`](https://golang.org/doc/install): The language Tektoncd-pipeline-operator is
   built in
1. [`git`](https://help.github.com/articles/set-up-git/): For source control
1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/): For
   interacting with your kube cluster
1. operator-sdk: https://github.com/operator-framework/operator-sdk
1. [ko](Option)(https://github.com/google/ko): Build and deploy Go applications on Kubernetes (optional)

# Details
The project implements `Controller/reconciler` based on `operator-sdk` and enhance it to use `ko` as build/deploy tool.

# Installation
1. Git clone the repo.
2. ko apply -f ./deploy
3. kubectl apply -f ./samples/tekton_v1alpha1_resourcewatcher_cr.yaml

