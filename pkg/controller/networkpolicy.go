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

	networkingv1 "k8s.io/api/networking/v1"
)

// NewNetworkPolicyController
func NewNetworkPolicyController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watcher.NewNetworkPolicyListWatch(clientset)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    AddNftlbPolicies,
		DeleteFunc: DeleteNftlbPolicies,
		UpdateFunc: UpdateNftlbPolicies,
	}

	_, controller := cache.NewInformer(
		listWatch,
		&networkingv1.NetworkPolicy{},
		0,
		eventHandler,
	)

	return controller
}

// AddNftlbPolicies
func AddNftlbPolicies(obj interface{}) {
	// Parse this Network policy object as a Policies struct
	policies, err := parser.NetworkPolicyAsPolicies(obj.(*networkingv1.NetworkPolicy))
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}

	// Parse Policies struct as a JSON string
	policiesJSON, err := parser.StructAsJSON(policies)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}

	// Fill the request data for policies
	policiesRequestData := &types.RequestData{
		Method: "POST",
		Path:   "policies",
		Body:   strings.NewReader(policiesJSON),
	}

	// Get the response from that request
	policiesResponse, err := http.Send(policiesRequestData)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}
	go log.WriteLog(types.StandardLog, string(policiesResponse))

	// Parse this Network policy object as a Farms struct
	farms, err := parser.NetworkPolicyAsFarms(obj.(*networkingv1.NetworkPolicy))
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}

	// Parse Policies struct as a JSON string
	farmsJSON, err := parser.StructAsJSON(farms)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}

	// Fill the request data for applying those policies to farms
	farmsRequestData := &types.RequestData{
		Method: "POST",
		Path:   "farms",
		Body:   strings.NewReader(farmsJSON),
	}

	// Get the response from that request
	farmsResponse, err := http.Send(farmsRequestData)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}
	go log.WriteLog(types.StandardLog, string(farmsResponse))
}

// DeleteNftlbPolicies
func DeleteNftlbPolicies(obj interface{}) {
	// Parse network policy names as a Policies struct
	policies := parser.NetworkPolicyNamesAsPolicies(obj.(*networkingv1.NetworkPolicy))

	// Parse Policies struct as a JSON string
	policiesJSON, err := parser.StructAsJSON(policies)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}

	// Fill the request data
	policiesRequestData := &types.RequestData{
		Method: "DELETE",
		Path:   "policies",
		Body:   strings.NewReader(policiesJSON),
	}

	// Get the response from that request
	policiesResponse, err := http.Send(policiesRequestData)
	if err != nil {
		go log.WriteLog(types.ErrorLog, err.Error())
		return
	}
	go log.WriteLog(types.StandardLog, string(policiesResponse))
}

// UpdateNftlbPolicies
func UpdateNftlbPolicies(oldObj, newObj interface{}) {
	DeleteNftlbPolicies(oldObj)
	AddNftlbPolicies(newObj)
}
