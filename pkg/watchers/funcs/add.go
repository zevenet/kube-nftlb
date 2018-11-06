package funcs

import (
	"fmt"
	"strings"

	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	json "github.com/zevenet/kube-nftlb/pkg/json"
	request "github.com/zevenet/kube-nftlb/pkg/request"
	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// CreateNftlbObject creates any object from nftlb (farm or backend) given its Kubernetes resource name.
func CreateNftlbObject(resourceName string, obj interface{}) {
	switch resourceName {
	case "Service":
		createNftlbFarm(obj)
	case "Endpoint":
		createNftlbBackends(obj)
	default:
		err := fmt.Sprintf("Resource not recognised: %s", resourceName)
		panic(err)
	}
}

// createNftlbFarm creates any nftlb farm given its name and the Service object.
func createNftlbFarm(service interface{}) {
	// Translates the Service object into a JSONnftlb struct
	JSONnftlb := json.GetJSONnftlbFromService(service)
	// Translates that struct into a JSON string
	farmJSON := json.DecodeJSON(JSONnftlb)
	// Makes the request
	createNftlbRequest(farmJSON)
}

// createNftlbFarm creates backends for any farm given its name and the Endpoints object.
func createNftlbBackends(endpoints interface{}) {
	// Translates the Endpoints object into a JSONnftlb struct
	JSONnftlb := json.GetJSONnftlbFromEndpoints(endpoints)
	// Translates that struct into a JSON string
	backendsJSON := json.DecodeJSON(JSONnftlb)
	// Makes the request
	createNftlbRequest(backendsJSON)
}

func createNftlbRequest(json string) {
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
