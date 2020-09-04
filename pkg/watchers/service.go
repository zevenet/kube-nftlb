package watchers

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

// NewServiceListWatch makes a ListWatch of every Service in the cluster.
func NewServiceListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	return cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), // REST interface
		"services",                      // Resource to watch for
		corev1.NamespaceAll,             // Resource can be found in ALL namespaces
		fields.Everything(),             // Get ALL fields from requested resource
	)
}
