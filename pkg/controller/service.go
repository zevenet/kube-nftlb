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
	svc := obj.(*corev1.Service)

	// Parse this Service struct as a Farms struct
	farms := parser.ServiceAsFarms(svc)

	// Don't accept empty farms
	if farms.Farms == nil || len(farms.Farms) == 0 {
		go log.WriteLog(types.DetailedLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\nFarms struct is empty", svc.Name))
		return
	}

	// Parse Farms struct as a parser string
	farmsJSON, err := parser.StructAsJSON(farms)
	if err != nil {
		go log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
		return
	}

	// Fill the request data for farms
	requestData := &types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(farmsJSON),
	}

	// Get the response from that request
	response, err := http.Send(requestData)
	if err != nil {
		go log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
		return
	}
	go log.WriteLog(types.StandardLog, fmt.Sprintf("AddNftlbFarms: Service name: %s\n%s", svc.Name, string(response)))
}

// DeleteNftlbFarms
func DeleteNftlbFarms(obj interface{}) {
	svc := obj.(*corev1.Service)

	// Make channel where farm path will arrive
	farmPathsChan := make(chan string, 1)

	// Handle shared channel
	go parser.DeleteMaxConnsService(svc)
	go parser.DeleteServiceFarms(svc, farmPathsChan)

	for farmPath := range farmPathsChan {
		// Fill the request data
		requestData := &types.RequestData{
			Method: "DELETE",
			Path:   farmPath,
		}

		// Get the response from that request
		if response, err := http.Send(requestData); err != nil {
			go log.WriteLog(types.ErrorLog, fmt.Sprintf("DeleteNftlbFarms: Service name: %s\n%s", svc.Name, err.Error()))
		} else {
			go log.WriteLog(types.StandardLog, fmt.Sprintf("DeleteNftlbFarms: Service name: %s\n%s", svc.Name, string(response)))
		}
	}
}

// UpdateNftlbFarms
func UpdateNftlbFarms(oldObj, newObj interface{}) {
	AddNftlbFarms(newObj.(*corev1.Service))
}
