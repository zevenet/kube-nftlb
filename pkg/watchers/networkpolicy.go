package watchers

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

// NewNetworkPolicyListWatch makes a ListWatch of every Network policy in the cluster.
func NewNetworkPolicyListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	return cache.NewListWatchFromClient(
		clientset.NetworkingV1().RESTClient(), // REST interface
		"networkpolicies",                     // Resource to watch for
		corev1.NamespaceAll,                   // Resource can be found in ALL namespaces
		fields.Everything(),                   // Get ALL fields from requested resource
	)
}
