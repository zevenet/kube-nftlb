package watchers

import (
	"fmt"

	logs "github.com/zevenet/kube-nftlb/pkg/logs"
	funcs "github.com/zevenet/kube-nftlb/pkg/watchers/funcs"
	v1 "k8s.io/api/core/v1"
	fields "k8s.io/apimachinery/pkg/fields"
	runtime "k8s.io/apimachinery/pkg/runtime"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

// getListWatch makes a ListWatch of every resource in the cluster.
func getListWatch(clientset *kubernetes.Clientset, resource string) *cache.ListWatch {
	listwatch := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), // REST interface
		resource,                        // Resource to watch for
		v1.NamespaceAll,                 // Resource can be found in ALL namespaces
		fields.Everything(),             // Get ALL fields from requested resource
	)
	return listwatch
}

// getController returns a Controller based on listWatch.
// Exports every message into logChannel.
func getController(listWatch *cache.ListWatch, resourceStruct runtime.Object, resourceName string, logChannel chan string, clientset *kubernetes.Clientset) cache.Controller {
	_, controller := cache.NewInformer(
		listWatch,      // Resources to watch for
		resourceStruct, // Resource struct
		0,
		// Event handler: new, deleted or updated resource
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				switch tp := obj.(type) {
				case *v1.Service:
					funcs.CreateNftlbFarm(obj.(*v1.Service), clientset, logChannel)
				case *v1.Endpoints:
					funcs.CreateNftlbBackends(obj.(*v1.Endpoints), logChannel, clientset)
				default:
					err := fmt.Sprintf("Object not recognised of type %t", tp)
					panic(err)
				}
				levelLog := 3
				logChannel = logs.PrintLogChannelFuncGeneral(levelLog, "New", resourceName, obj, logChannel)
			},
			DeleteFunc: func(obj interface{}) {
				switch tp := obj.(type) {
				case *v1.Service:
					funcs.DeleteNftlbFarm(obj.(*v1.Service), logChannel)
				case *v1.Endpoints:
					funcs.DeleteNftlbBackends(obj.(*v1.Endpoints), logChannel)
				default:
					err := fmt.Sprintf("Object not recognised of type %t", tp)
					panic(err)
				}
				levelLog := 3
				logChannel = logs.PrintLogChannelFuncGeneral(levelLog, "nDeleted", resourceName, obj, logChannel)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				switch tp := oldObj.(type) {
				case *v1.Service:
					funcs.UpdateNftlbFarm(newObj.(*v1.Service), clientset, logChannel)
				case *v1.Endpoints:
					funcs.UpdateNftlbBackends(oldObj.(*v1.Endpoints), newObj.(*v1.Endpoints), logChannel, clientset)
				default:
					err := fmt.Sprintf("Object not recognised of type %t", tp)
					panic(err)
				}
				levelLog := 3
				logChannel = logs.PrintLogChannelFuncUpdate(levelLog, "\nUpdated %s:\n* BEFORE: %s\n* NOW: %s", resourceName, oldObj, newObj, logChannel)
			},
		},
	)
	return controller
}
