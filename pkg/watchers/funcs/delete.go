package funcs

import (
	"fmt"

	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	request "github.com/zevenet/kube-nftlb/pkg/request"
	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// DeleteNftlbObject deletes any object from nftlb (farm or backend) given its Kubernetes resource name.
func DeleteNftlbObject(resourceName string, obj interface{}) {
	switch resourceName {
	case "Service":
		deleteNftlbFarmFromService(obj)
	case "Endpoint":
		deleteNftlbFarmFromEndpoints(obj)
	default:
		err := fmt.Sprintf("Resource not recognised: %s", resourceName)
		panic(err)
	}
}

// deleteNftlbFarmFromService deletes any nftlb farm given a Service object.
func deleteNftlbFarmFromService(service interface{}) {
	serviceObj := service.(v1.Service)
	farmName := serviceObj.GetLabels()["app"]
	deleteNftlbFarm(farmName)
}

// deleteNftlbFarmFromEndpoints deletes any nftlb farm given a Endpoints object.
func deleteNftlbFarmFromEndpoints(endpoints interface{}) {
	endpointsObj := endpoints.(v1.Endpoints)
	farmName := endpointsObj.GetLabels()["app"]
	deleteNftlbFarm(farmName)
}

// deleteNftlbFarm deletes any nftlb farm given its name.
func deleteNftlbFarm(name string) {
	// Makes the farm path
	farmPath := fmt.Sprintf("/%s", name)
	// Makes the URL and its Header
	farmURL := defaults.SetNftlbURL(farmPath)
	nftlbKey := defaults.SetNftlbKey()
	// Fills the request
	rq := &types.Request{
		Header: nftlbKey,
		Action: types.DELETE,
		URL:    farmURL,
	}
	// Does the request
	resp := request.GetResponse(rq)
	// Shows the response
	fmt.Println(resp)
}
