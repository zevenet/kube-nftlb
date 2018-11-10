package funcs

import (
	"fmt"

	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	json "github.com/zevenet/kube-nftlb/pkg/json"
	request "github.com/zevenet/kube-nftlb/pkg/request"
	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// DeleteNftlbFarm deletes any nftlb farm given a Service object.
func DeleteNftlbFarm(service *v1.Service) {
	farmName := service.ObjectMeta.Name
	deleteNftlbRequest(farmName)
}

// DeleteNftlbBackends deletes all nftlb backends from a farm given a Endpoints object.
func DeleteNftlbBackends(endpoints *v1.Endpoints) {
	farmName := endpoints.ObjectMeta.Name
	for json.GetBackendID(farmName) > 0 {
		backendName := fmt.Sprintf("%s%d", farmName, json.GetBackendID(farmName))
		fullPath := fmt.Sprintf("%s/backends/%s", farmName, backendName)
		deleteNftlbRequest(fullPath)
		json.DecreaseBackendID(farmName)
	}
}

func deleteNftlbRequest(name string) {
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
