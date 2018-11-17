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
		// Makes the request
		response := updateNftlbRequest(farmJSON)
		// Prints info
		printUpdated("Farm", farmJSON, response)
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
		// Makes the request
		response := updateNftlbRequest(backendsJSON)
		// Prints info
		printUpdated("Backends", backendsJSON, response)
		// Deletes remaining old backends if any
		newNumberBackends := json.GetBackendID(farmName)
		for oldNumberBackends > newNumberBackends {
			// Decreases number of old backends
			oldNumberBackends--
			// Makes the full path for the request
			backendName := fmt.Sprintf("%s%d", farmName, oldNumberBackends)
			fullPath := fmt.Sprintf("%s/backends/%s", farmName, backendName)
			response := deleteNftlbRequest(fullPath)
			// Prints info
			printDeleted("Backend", farmName, backendName, response)
		}
	}
}

func updateNftlbRequest(json string) string {
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
	// Returns the response
	return request.GetResponse(rq)
}

func printUpdated(object string, json string, response string) {
	message := fmt.Sprintf("\nUpdated %s:\n%s\n%s", object, json, response)
	fmt.Println(message)
}
