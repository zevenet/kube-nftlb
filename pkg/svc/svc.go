package svc

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	fields "k8s.io/apimachinery/pkg/fields"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

// GetServiceListWatch makes a ListWatch of every Services in the cluster.
func GetServiceListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	listwatch := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), // REST interface
		string(v1.ResourceServices),     // Resource to watch for: "Service"
		v1.NamespaceAll,                 // Resource can be found in ALL namespaces
		fields.Everything(),             // Get ALL fields from requested resource
	)
	return listwatch
}

// GetServiceController returns a Controller based on listWatch.
// Exports every message into logChannel
func GetServiceController(listWatch *cache.ListWatch, logChannel chan string) cache.Controller {
	_, controller := cache.NewInformer(
		listWatch,     // Resources to watch for
		&v1.Service{}, // Resource struct
		0,
		// Event handler: new, deleted or updated resource
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				logChannel <- fmt.Sprintf("New Service: %s\n\n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				logChannel <- fmt.Sprintf("Deleted Service: %s\n\n", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				logChannel <- fmt.Sprintf("Updated Service:\nBEFORE: %s\nNOW: %s\n\n", oldObj, newObj)
			},
		},
	)
	return controller
}
