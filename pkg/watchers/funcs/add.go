package funcs

import (
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/http"
	"github.com/zevenet/kube-nftlb/pkg/json"
	"github.com/zevenet/kube-nftlb/pkg/logs"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/core/v1"
)

// CreateNftlbFarm creates any nftlb farm given a Service object.
func CreateNftlbFarm(service *v1.Service, clientset *kubernetes.Clientset) {
	if !json.Contains(http.BadNames, service.ObjectMeta.Name) {
		// Translates the Service object into a JSONnftlb struct
		JSONnftlb := json.GetJSONnftlbFromService(service, clientset)
		// Translates that struct into a JSON string
		farmJSON := json.DecodePrettyJSON(JSONnftlb)
		// Makes the request
		response := createNftlbRequest(farmJSON)
		// Prints info
		printNew("Farm", farmJSON, response)
	}
}

// CreateNftlbBackends creates backends for any farm given a Endpoints object.
func CreateNftlbBackends(endpoints *v1.Endpoints, clientset *kubernetes.Clientset) {
	if !json.Contains(http.BadNames, endpoints.ObjectMeta.Name) {
		// Translates the Endpoints object into a JSONnftlb struct
		JSONnftlb := json.GetJSONnftlbFromEndpoints(endpoints, clientset)
		// Translates that struct into a JSON string
		backendsJSON := json.DecodePrettyJSON(JSONnftlb)
		// Makes the request
		response := createNftlbRequest(backendsJSON)
		// Prints info
		printNew("Backends", backendsJSON, response)
	}
}

func createNftlbRequest(json string) string {
	// Fill the request data
	requestData := &types.RequestData{
		Method: "POST",
		Body:   strings.NewReader(json),
	}

	// Get the response from that request
	response, err := http.Send(requestData)
	if err != nil {
		panic(err)
	}

	return string(response)
}

func printNew(object string, json string, response string) {
	levelLog := 0
	message := fmt.Sprintf("\nNew %s:\n%s\n%s", object, json, response)
	logs.WriteLog(levelLog, message)
}
