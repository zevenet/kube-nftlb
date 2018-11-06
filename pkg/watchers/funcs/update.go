package funcs

import (
	"fmt"
	"strings"

	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	json "github.com/zevenet/kube-nftlb/pkg/json"
	request "github.com/zevenet/kube-nftlb/pkg/request"
	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// UpdateNftlbObject updates any nftlb farm or backend given both (updated and old) objects.
func UpdateNftlbObject(resourceName string, oldObj, newObj interface{}) {
	switch resourceName {
	case "Service":
		updateNftlbFarm(newObj)
	case "Endpoint":
		updateNftlbBackends(oldObj, newObj)
	default:
		err := fmt.Sprintf("Resource not recognised: %s", resourceName)
		panic(err)
	}
}

// updateNftlbFarm updates any nftlb farm given its name and the Service object.
func updateNftlbFarm(newSvc interface{}) {
	// Translates the updated Service object into a JSONnftlb struct
	newJSONnftlb := json.GetJSONnftlbFromService(newSvc)
	// Translates that struct into a JSON string
	farmJSON := json.DecodeJSON(newJSONnftlb)
	// Makes the request
	updateNftlbRequest(farmJSON)
}

// updateNftlbFarm updates backends for any farm given its name and the Endpoints objects.
func updateNftlbBackends(oldEP, newEP interface{}) {
	// Translates the Endpoints objects into JSONnftlb structs
	oldJSONnftlb := json.GetJSONnftlbFromEndpoints(oldEP)
	newJSONnftlb := json.GetJSONnftlbFromEndpoints(newEP)
	/**
	* TODO: Compare both objects to know which endpoints have changed
	 */
	// Translates that struct into a JSON string
	backendsJSON := json.DecodeJSON(newJSONnftlb)
	// Makes the request
	updateNftlbRequest(backendsJSON)
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
