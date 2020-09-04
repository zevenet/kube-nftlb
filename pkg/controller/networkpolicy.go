package controller

import (
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/http"
	"github.com/zevenet/kube-nftlb/pkg/json"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"github.com/zevenet/kube-nftlb/pkg/watchers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	networkingv1 "k8s.io/api/networking/v1"
)

// NewNetworkPolicyController
func NewNetworkPolicyController(clientset *kubernetes.Clientset) cache.Controller {
	listWatch := watchers.NewNetworkPolicyListWatch(clientset)

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
	policies, err := json.ParseNetworkPolicyAsPolicies(obj.(*networkingv1.NetworkPolicy))
	if err != nil {
		// Log error if it couldn't be parsed
		return
	}

	// Parse Policies struct as a JSON string
	policiesJSON, err := json.ParseStruct(policies)
	if err != nil {
		// Log error if it couldn't be parsed
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
		// Log error if the request failed
		return
	}

	// Log response
	fmt.Println(policiesResponse)

	// Parse this Network policy object as a Farms struct
	farms, err := json.ParseNetworkPolicyAsFarms(obj.(*networkingv1.NetworkPolicy))
	if err != nil {
		// Log error if it couldn't be parsed
		return
	}

	// Parse Policies struct as a JSON string
	farmsJSON, err := json.ParseStruct(farms)
	if err != nil {
		// Log error if it couldn't be parsed
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
		// Log error if the request failed
		return
	}

	// Log response
	fmt.Println(farmsResponse)
}

// DeleteNftlbPolicies
func DeleteNftlbPolicies(obj interface{}) {
	// Parse network policy names as a Policies struct
	policies := json.ParseNetworkPolicyNamesAsPolicies(obj.(*networkingv1.NetworkPolicy))

	// Parse Policies struct as a JSON string
	policiesJSON, err := json.ParseStruct(policies)
	if err != nil {
		// Log error if it couldn't be parsed
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
		// Log error if the request failed
		return
	}

	// Log response
	fmt.Println(policiesResponse)
}

// UpdateNftlbPolicies
func UpdateNftlbPolicies(oldObj, newObj interface{}) {
	DeleteNftlbPolicies(oldObj)
	AddNftlbPolicies(newObj)
}
