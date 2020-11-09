package controller

import (
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/http"
	"github.com/zevenet/kube-nftlb/pkg/log"
	"github.com/zevenet/kube-nftlb/pkg/metrics"
	"github.com/zevenet/kube-nftlb/pkg/parser"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"github.com/zevenet/kube-nftlb/pkg/watcher"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

// NewEndpointsController
func NewEndpointsController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watcher.NewEndpointListWatch(clientset)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    AddNftlbBackends,
		DeleteFunc: DeleteNftlbBackends,
		UpdateFunc: UpdateNftlbBackends,
	}

	_, controller := cache.NewInformer(
		listWatch,
		&corev1.Endpoints{},
		0,
		eventHandler,
	)

	return controller
}

// AddNftlbBackends
func AddNftlbBackends(obj interface{}) {
	ep := obj.(*corev1.Endpoints)

	// Parse this Endpoints struct as a Nftlb struct
	data := parser.EndpointsAsNftlb(ep)

	if len(data.Farms) == 0 {
		// Reject object without farms
		log.WriteLog(types.DetailedLog, fmt.Sprintf("AddNftlbFarms: Endpoints name: %s\nEmpty Farms slice", ep.Name))
		return
	} else if len(data.Farms[0].Backends) == 0 {
		// Reject farm without backends
		log.WriteLog(types.DetailedLog, fmt.Sprintf("AddNftlbFarms: Endpoints name: %s\nFarms[0] has no backends", ep.Name))
		return
	}

	// Parse Nftlb struct as a JSON string
	nftlbJSON, err := parser.NftlbAsJSON(data)
	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, err.Error()))
		return
	}
	log.WriteLog(types.StandardLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, nftlbJSON))

	metrics.EndpointsChangesPending.Inc()
	metrics.EndpointsChangesTotal.Inc()
	// Get the response from that request
	response, err := http.Send(&types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(nftlbJSON),
	})
	metrics.EndpointsChangesPending.Dec()

	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, err.Error()))
		return
	}

	log.WriteLog(types.StandardLog, string(response))
}

// DeleteNftlbBackends
func DeleteNftlbBackends(obj interface{}) {
	metrics.EndpointsChangesTotal.Inc()
	ep := obj.(*corev1.Endpoints)
	pathsChan := make(chan string)

	go func() {
		for path := range pathsChan {
			// Get the response from that request
			if response, err := http.Send(&types.RequestData{
				Method: "DELETE",
				Path:   path,
			}); err != nil {
				log.WriteLog(types.ErrorLog, fmt.Sprintf("DeleteNftlbBackends: Endpoints name: %s, path: %s\n%s", ep.Name, path, err.Error()))
			} else {
				log.WriteLog(types.StandardLog, fmt.Sprintf("DeleteNftlbBackends: Endpoints name: %s, path: %s\n%s", ep.Name, path, string(response)))
			}
		}
	}()

	parser.EndpointsAsPaths(ep, pathsChan)
}

// UpdateNftlbBackends
func UpdateNftlbBackends(oldObj, newObj interface{}) {
	DeleteNftlbBackends(oldObj)
	AddNftlbBackends(newObj)
}
