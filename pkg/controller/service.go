package controller

import (
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

// NewServiceController
func NewServiceController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watcher.NewServiceListWatch(clientset)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    AddNftlbFarms,
		DeleteFunc: DeleteNftlbFarms,
		UpdateFunc: UpdateNftlbFarms,
	}

	_, controller := cache.NewInformer(
		listWatch,
		&corev1.Service{},
		0,
		eventHandler,
	)

	return controller
}

// AddNftlbFarms
func AddNftlbFarms(obj interface{}) {
	// Parse this Service struct as a Farms struct
	farms := parser.ServiceAsFarms(obj.(*corev1.Service))

	// Don't accept empty farms
	if farms.Farms == nil || len(farms.Farms) == 0 {
		return
	}

	// Parse Farms struct as a parser string
	farmsJSON, err := parser.StructAsJSON(farms)
	if err != nil {
		// Log error if it couldn't be parsed
		return
	}

	go log.WriteLog(0, farmsJSON)

	// Fill the request data for farms
	requestData := &types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(farmsJSON),
	}

	// Get the response from that request
	if _, err := http.Send(requestData); err != nil {
		// Log error if the request failed
		return
	}
}

// DeleteNftlbFarms
func DeleteNftlbFarms(obj interface{}) {
	// Make channel where farm path will arrive
	farmPathsChan := make(chan string, 1)

	// Handle shared channel
	go parser.DeleteMaxConnsService(obj.(*corev1.Service))
	go parser.DeleteServiceFarms(obj.(*corev1.Service), farmPathsChan)

	for farmPath := range farmPathsChan {
		// Fill the request data
		requestData := &types.RequestData{
			Method: "DELETE",
			Path:   farmPath,
		}

		// Get the response from that request
		if _, err := http.Send(requestData); err != nil {
			// Log error if the request failed
		}
	}
}

// UpdateNftlbFarms
func UpdateNftlbFarms(oldObj, newObj interface{}) {
	AddNftlbFarms(newObj.(*corev1.Service))
}
