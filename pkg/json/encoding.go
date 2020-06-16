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
func GetJSONnftlbFromService(service *v1.Service) types.JSONnftlb {
	farmName := service.ObjectMeta.Name
	// Extracts ports and protocols as strings
	protocolsSlice := map[string]int{"TCP": 0, "UDP": 0}
	var portsSlice []string
	// For every service port:
	for _, port := range service.Spec.Ports {
		// Gets the port as string
		portString := fmt.Sprint(port.Port)
		// Increase number of protocol ocurrences
		protocolsSlice[string(port.Protocol)]++
		// Don't duplicate if it exists already
		if !Contains(portsSlice, portString) {
			portsSlice = append(portsSlice, portString)
		}
	}
	// Ports slice -> port(s) string
	ports := strings.Join(portsSlice, ", ")
	// Gets the protocol
	protocol := ""
	if protocolsSlice["TCP"] > 0 && protocolsSlice["UDP"] > 0 {
		protocol = "all"
	} else if protocolsSlice["TCP"] > 0 {
		protocol = "tcp"
	} else if protocolsSlice["UDP"] > 0 {
		protocol = "udp"
	}
	// Gets persistence and Stickiness timeout in seconds
        persistence := ""
        if service.Spec.SessionAffinity == "ClientIP"{
                persistence = "srcip"
        }else if service.Spec.SessionAffinity == "None"{
                persistence = "none"
        }
        persistence_ttl := ""
        if service.Spec.SessionAffinityConfig != nil{
                if service.Spec.SessionAffinityConfig.ClientIP != nil{
                        if service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds != nil{
				// Value between 0 and 86400 seconds (1 day max)
                                persistence_ttl = fmt.Sprint(*(service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds))
                        }
                }
        }

	// Fills the farm
	var farm = types.Farms{
		types.Farm{
			Name:         farmName,
			Family:       "ipv4",
			VirtualAddr:  service.Spec.ClusterIP,
			VirtualPorts: ports,
			Mode:         "snat",
			Protocol:     protocol,
			State:        "up",
			Intraconnect: "on",
			Persistence:  fmt.Sprint(persistence),
                        PersistTTL:   fmt.Sprint(persistence_ttl),
			Backends:     types.Backends{},
		},
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farm,
	}
}

// GetJSONnftlbFromEndpoints returns a JSONnftlb struct filled with any Endpoints data.
func GetJSONnftlbFromEndpoints(endpoints *v1.Endpoints) types.JSONnftlb {
	farmName := endpoints.ObjectMeta.Name
	// Extracts individual addresses
	var addrSlice []string
	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			addrSlice = append(addrSlice, address.IP)
		}
	}
	// Initializes farm/backends ID
	CreateFarmID(farmName)
	// Fills backends
	var backends types.Backends
	for _, address := range addrSlice {
		// Gets backend ID
		backendID := GetBackendID(farmName)

		// Create backend for each port
		for _, endpoint2 := range endpoints.Subsets {
                	for _, port := range endpoint2.Ports {
				var backend = types.Backend{
				Name:   fmt.Sprintf("%s%d", farmName, backendID),
				IPAddr: address,
				State:  "up",
				Port: fmt.Sprint(port.Port),
                	}
			// Appends backend
			backends = append(backends, backend)
                }
		// Increases backend ID
		IncreaseBackendID(farmName)
		}
	}
	// Fills the farm
	var farm = types.Farms{
		types.Farm{
			Name:     farmName,
			Backends: backends,
		},
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farm,
	}
}

// Contains returns true when "str" string is in "sl" slice.
func Contains(sl []string, str string) bool {
	for _, value := range sl {
		if value == str {
			return true
		}
	}
	return false
}
