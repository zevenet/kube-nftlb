package watchers

import (
	v1 "k8s.io/api/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

var (
	resourceNameEP   = "endpoints"
	resourceStructEP = v1.Endpoints{}
)

// GetEndpointListWatch makes a ListWatch of every Endpoint in the cluster.
func GetEndpointListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	return getListWatch(clientset, resourceNameEP)
}

// GetEndpointController returns a Controller based on listWatch.
// Exports every message into logChannel.
func GetEndpointController(listWatch *cache.ListWatch, logChannel chan string, clientset *kubernetes.Clientset) cache.Controller {
	return getController(listWatch, &resourceStructEP, "Endpoint", logChannel, clientset)
}
