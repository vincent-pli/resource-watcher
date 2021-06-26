module github.com/vincent-pli/resource-watcher

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.12.7 // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.1.1 // indirect
	github.com/Azure/go-autorest v12.2.0+incompatible
	github.com/cloudevents/sdk-go v0.9.2 // indirect
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.0
	github.com/operator-framework/operator-sdk v0.11.1-0.20191012024916-f419ad3f3dc5
	github.com/spf13/pflag v1.0.3
	google.golang.org/grpc v1.22.1 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208
	knative.dev/eventing v0.9.0
	knative.dev/pkg v0.0.0-20191016060315-3f11504864ae
	sigs.k8s.io/controller-runtime v0.2.0
)

// Pinned to kubernetes-1.13.4
replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	// Pinned to v2.9.2 (kubernetes-1.13.1) so https://proxy.golang.org can
	// resolve it correctly.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20190424153033-d3245f150225
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.0
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.0
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.11.0

replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822

go 1.16
