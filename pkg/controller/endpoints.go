package controller

import (
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/http"
	"github.com/zevenet/kube-nftlb/pkg/json"
	"github.com/zevenet/kube-nftlb/pkg/logs"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"github.com/zevenet/kube-nftlb/pkg/watchers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

// NewEndpointsController
func NewEndpointsController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watchers.NewEndpointListWatch(clientset)

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
	// Parse this Service struct as a Farms struct
	farms := json.ParseEndpointsAsFarms(obj.(*corev1.Endpoints))

	// Don't accept empty farms
	if farms.Farms == nil || len(farms.Farms) == 0 {
		return
	}

	// Parse Farms struct as a JSON string
	farmsJSON, err := json.ParseStruct(farms)
	if err != nil {
		// Log error if it couldn't be parsed
		return
	}

	go logs.WriteLog(0, farmsJSON)

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

// DeleteNftlbBackends
func DeleteNftlbBackends(obj interface{}) {
	backendsChan := make(chan string, 1)

	go json.DeleteEndpointsBackends(obj.(*corev1.Endpoints), backendsChan)

	for backendPath := range backendsChan {
		// Fills the request data
		requestData := &types.RequestData{
			Method: "DELETE",
			Path:   backendPath,
		}

		// Get the response from that request
		response, err := http.Send(requestData)
		if err != nil {
			// Log error
			continue
		}

		// Log response
		logs.WriteLog(0, string(response))
	}
}

// UpdateNftlbBackends
func UpdateNftlbBackends(oldObj, newObj interface{}) {
	DeleteNftlbBackends(oldObj.(*corev1.Endpoints))
	AddNftlbBackends(newObj.(*corev1.Endpoints))
}
