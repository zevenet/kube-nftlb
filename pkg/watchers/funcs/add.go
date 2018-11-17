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

// CreateNftlbFarm creates any nftlb farm given a Service object.
func CreateNftlbFarm(service *v1.Service) {
	if !json.Contains(request.BadNames, service.ObjectMeta.Name) {
		// Translates the Service object into a JSONnftlb struct
		JSONnftlb := json.GetJSONnftlbFromService(service)
		// Translates that struct into a JSON string
		farmJSON := json.DecodePrettyJSON(JSONnftlb)
		// Makes the request
		response := createNftlbRequest(farmJSON)
		// Prints info
		printNew("Farm", farmJSON, response)
	}
}

// CreateNftlbBackends creates backends for any farm given a Endpoints object.
func CreateNftlbBackends(endpoints *v1.Endpoints) {
	if !json.Contains(request.BadNames, endpoints.ObjectMeta.Name) {
		// Translates the Endpoints object into a JSONnftlb struct
		JSONnftlb := json.GetJSONnftlbFromEndpoints(endpoints)
		// Translates that struct into a JSON string
		backendsJSON := json.DecodePrettyJSON(JSONnftlb)
		// Makes the request
		response := createNftlbRequest(backendsJSON)
		// Prints info
		printNew("Backends", backendsJSON, response)
	}
}

func createNftlbRequest(json string) string {
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

func printNew(object string, json string, response string) {
	message := fmt.Sprintf("\nNew %s:\n%s\n%s", object, json, response)
	fmt.Println(message)
}
