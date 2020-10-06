package controller

import (
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/http"
	"github.com/zevenet/kube-nftlb/pkg/log"
	"github.com/zevenet/kube-nftlb/pkg/parser"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"github.com/zevenet/kube-nftlb/pkg/watcher"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

// NewServiceController returns a k8s controller with a Service resource watcher, and runs different functions based on the
// event type that the watcher notifies.
func NewServiceController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watcher.NewServiceListWatch(clientset)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    AddNftlbFarm,
		DeleteFunc: DeleteNftlbFarm,
		UpdateFunc: UpdateNftlbFarm,
	}

	_, controller := cache.NewInformer(
		listWatch,
		&corev1.Service{},
		0,
		eventHandler,
	)

	return controller
}

// AddNftlbFarm takes in a Service object (k8s) and creates a farm with addresses (nftlb).
func AddNftlbFarm(obj interface{}) {
	svc := obj.(*corev1.Service)

	// Reject an invalid Service
	if svc.Spec.ClusterIP == "" {
		log.WriteLog(types.DetailedLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\nInvalid Service, ClusterIP should not be empty", svc.Name))
		return
	}

	// Parse Service as a Nftlb struct
	data := parser.ServiceAsNftlb(svc)

	// Parse Nftlb struct as JSON
	nftlbJSON, err := parser.NftlbAsJSON(data)
	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
		return
	}

	// Send that JSON data to nftlb
	response, err := http.Send(&types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(nftlbJSON),
	})
	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
		return
	}

	// Read the response
	log.WriteLog(types.StandardLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, string(response)))
}

// DeleteNftlbFarm takes in a Service object (k8s) and deletes the farm related to the service and its addresses (nftlb).
func DeleteNftlbFarm(obj interface{}) {
	svc := obj.(*corev1.Service)

	// Make channel where paths will come through
	pathChan := make(chan string)

	go func() {
		for path := range pathChan {
			// Get the response from that request
			if response, err := http.Send(&types.RequestData{
				Method: "DELETE",
				Path:   path,
			}); err != nil {
				log.WriteLog(types.ErrorLog, fmt.Sprintf("DeleteNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
			} else {
				log.WriteLog(types.StandardLog, fmt.Sprintf("DeleteNftlbFarms: Service name: %s\n%s", svc.Name, string(response)))
			}
		}
	}()

	// Read paths and send them through the channel
	parser.ServiceAsPaths(svc, pathChan)
}

// UpdateNftlbFarm takes in two Services (both are the same, but one it's before the update and the other it's updated)
// and applies the changes from the updated Service.
func UpdateNftlbFarm(oldObj, newObj interface{}) {
	AddNftlbFarm(newObj)
}
