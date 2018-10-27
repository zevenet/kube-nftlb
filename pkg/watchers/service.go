package watchers

import (
	v1 "k8s.io/api/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

var (
	resourceNameSvc   = string(v1.ResourceServices)
	resourceStructSvc = v1.Service{}
)

// GetServiceListWatch makes a ListWatch of every Service in the cluster.
func GetServiceListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	return getListWatch(clientset, resourceNameSvc)
}

// GetServiceController returns a Controller based on listWatch.
// Exports every message into logChannel.
func GetServiceController(listWatch *cache.ListWatch, logChannel chan string) cache.Controller {
	return getController(listWatch, &resourceStructSvc, "Service", logChannel)
}
