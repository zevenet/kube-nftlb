package funcs

import (
	"fmt"
	"strings"

	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	json "github.com/zevenet/kube-nftlb/pkg/json"
	request "github.com/zevenet/kube-nftlb/pkg/request"
	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// UpdateNftlbFarm updates any nftlb farm given a Service object.
func UpdateNftlbFarm(newSvc *v1.Service) {
	if !json.Contains(request.BadNames, newSvc.ObjectMeta.Name) {
		// Translates the updated Service object into a JSONnftlb struct
		newJSONnftlb := json.GetJSONnftlbFromService(newSvc)
		// Translates that struct into a JSON string
		farmJSON := json.DecodePrettyJSON(newJSONnftlb)
		// Logs JSON
		fmt.Println("\nUpdated Service:")
		fmt.Println(farmJSON)
		// Makes the request
		updateNftlbRequest(farmJSON)
	}
}

// UpdateNftlbBackends updates backends for any farm given a Endpoints object.
func UpdateNftlbBackends(oldEP, newEP *v1.Endpoints) {
	if !json.Contains(request.BadNames, newEP.ObjectMeta.Name) {
		// Gets the farm name and number of backends for later
		farmName := oldEP.ObjectMeta.Name
		oldNumberBackends := json.GetBackendID(farmName)
		// Translates the Endpoints objects into JSONnftlb structs
		newJSONnftlb := json.GetJSONnftlbFromEndpoints(newEP)
		// Translates the struct into a JSON string
		backendsJSON := json.DecodePrettyJSON(newJSONnftlb)
		// Logs JSON
		fmt.Println("\nUpdated Endpoints:")
		fmt.Println(backendsJSON)
		// Makes the request
		updateNftlbRequest(backendsJSON)
		// Deletes remaining old backends if any
		newNumberBackends := json.GetBackendID(farmName)
		for oldNumberBackends > newNumberBackends {
			oldNumberBackends--
			backendName := fmt.Sprintf("%s%d", farmName, oldNumberBackends)
			fullPath := fmt.Sprintf("%s/backends/%s", farmName, backendName)
			deleteNftlbRequest(fullPath)
		}
	}
}

func updateNftlbRequest(json string) {
	// Makes the URL and its Header
	farmURL := defaults.SetNftlbURL("")
	nftlbKey := defaults.SetNftlbKey()
	// Fills the request
	rq := &types.Request{
		Header:  nftlbKey,
		Action:  types.POST,
		URL:     farmURL,
		Payload: strings.NewReader(json),
	}
	// Does the request
	resp := request.GetResponse(rq)
	// Shows the response
	fmt.Println(resp)
}
