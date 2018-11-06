package json

import (
	"encoding/json"
	"fmt"
	"strings"

	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// EncodeJSON returns a JSONnftlb struct with its fields filled with the JSON data.
func EncodeJSON(stringJSON string) types.JSONnftlb {
	var encodedJSON types.JSONnftlb
	if err := json.Unmarshal([]byte(stringJSON), &encodedJSON); err != nil {
		panic(err.Error())
	}
	return encodedJSON
}

// GetJSONnftlbFromService returns a JSONnftlb struct filled with any Service data.
func GetJSONnftlbFromService(serviceObj interface{}) types.JSONnftlb {
	service := serviceObj.(v1.Service)
	nameApp := service.GetLabels()["app"]
	// Extracts ports as strings
	var portsSlice []string
	for _, port := range service.Spec.Ports {
		portsSlice = append(portsSlice, port.String())
	}
	ports := strings.Join(portsSlice, ", ")
	// Fills the farm
	var farm = types.Farms{
		types.Farm{
			Name:         nameApp,
			Family:       "ipv4",
			VirtualAddr:  service.Spec.ClusterIP,
			VirtualPorts: ports,
			Mode:         "snat",
			Backends:     types.Backends{},
		},
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farm,
	}
}

// GetJSONnftlbFromEndpoints returns a JSONnftlb struct filled with any Endpoints data.
func GetJSONnftlbFromEndpoints(endpointsObj interface{}) types.JSONnftlb {
	endpoints := endpointsObj.(v1.Endpoints)
	nameApp := endpoints.GetLabels()["app"]
	// Extracts individual addresses
	var addrSlice []string
	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			addrSlice = append(addrSlice, address.String())
		}
	}
	// Initializes farm/backends ID
	CreateFarmID(nameApp)
	// Fills backends
	var backends types.Backends
	for _, address := range addrSlice {
		// Gets backend ID
		backendID := GetBackendID(nameApp)
		// Fills backend
		var backend = types.Backend{
			Name:   fmt.Sprintf("%s%d", nameApp, backendID),
			IPAddr: address,
		}
		// Appends backend
		backends = append(backends, backend)
		// Increases backend ID
		IncreaseBackendID(nameApp)
	}
	// Fills the farm
	var farm = types.Farms{
		types.Farm{
			Name:     nameApp,
			Backends: backends,
		},
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farm,
	}
}
