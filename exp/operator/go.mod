module sigs.k8s.io/cluster-api/exp/operator

go 1.16

require (
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/component-base v0.22.2
	k8s.io/klog/v2 v2.9.0
	k8s.io/kubectl v0.22.2
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/cluster-api v1.0.1-0.20211014163720-d2a110e53f48
	sigs.k8s.io/controller-runtime v0.10.2
)
