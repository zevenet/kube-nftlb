package controller

// TODO Adapt Endpoints to new Addresses nftlb object

/*
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

	// Parse this Service struct as a Farms struct
	farms := parser.EndpointsAsFarms(ep)

	// Don't accept empty farms
	if farms.Farms == nil || len(farms.Farms) == 0 {
		log.WriteLog(types.DetailedLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\nFarms struct is empty", ep.Name))
		return
	}

	// Parse Farms struct as a JSON string
	farmsJSON, err := parser.StructAsJSON(farms)
	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, err.Error()))
		return
	}
	log.WriteLog(types.StandardLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, farmsJSON))

	// Fill the request data for farms
	requestData := &types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(farmsJSON),
	}

	// Get the response from that request
	response, err := http.Send(requestData)
	if err != nil {
		log.WriteLog(types.ErrorLog, fmt.Sprintf("AddNftlbBackends: Endpoints name: %s\n%s", ep.Name, err.Error()))
		return
	}

	log.WriteLog(types.StandardLog, string(response))
}

// DeleteNftlbBackends
func DeleteNftlbBackends(obj interface{}) {
	ep := obj.(*corev1.Endpoints)
	backendsChan := make(chan string, 1)

	go parser.DeleteEndpointsBackends(ep, backendsChan)

	for backendPath := range backendsChan {
		// Fills the request data
		requestData := &types.RequestData{
			Method: "DELETE",
			Path:   backendPath,
		}

		// Get the response from that request
		if response, err := http.Send(requestData); err != nil {
			go log.WriteLog(types.ErrorLog, fmt.Sprintf("DeleteNftlbBackends: Endpoints name: %s, backend path: %s\n%s", ep.Name, backendPath, err.Error()))
		} else {
			go log.WriteLog(types.StandardLog, fmt.Sprintf("DeleteNftlbBackends: Endpoints name: %s, backend path: %s\n%s", ep.Name, backendPath, string(response)))
		}
	}
}

// UpdateNftlbBackends
func UpdateNftlbBackends(oldObj, newObj interface{}) {
	DeleteNftlbBackends(oldObj)
	AddNftlbBackends(newObj)
}
*/
